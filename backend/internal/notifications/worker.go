package notifications

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	workerBatchSize = 25
	maxAttempts     = 5
)

type Worker struct {
	repo   WorkerRepository
	sender Sender
	logger *log.Logger
}

type WorkerRepository interface {
	DispatchBatch(ctx context.Context, limit int) (int, error)
	ClaimDeliveries(ctx context.Context, limit int) ([]Delivery, error)
	MarkSent(ctx context.Context, deliveryID int64) error
	MarkRetry(ctx context.Context, retry Retry) error
	MarkFailed(ctx context.Context, deliveryID int64, message string) error
	DeactivateSubscription(ctx context.Context, subscriptionID int64) error
}

func NewWorker(repo WorkerRepository, sender Sender, logger *log.Logger) *Worker {
	return &Worker{repo: repo, sender: sender, logger: logger}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	w.process(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.process(ctx)
		}
	}
}

func (w *Worker) process(ctx context.Context) {
	if _, err := w.repo.DispatchBatch(ctx, workerBatchSize); err != nil {
		w.logger.Printf("notification dispatcher: %v", err)
		return
	}

	deliveries, err := w.repo.ClaimDeliveries(ctx, workerBatchSize)
	if err != nil {
		w.logger.Printf("notification sender: %v", err)
		return
	}

	for _, delivery := range deliveries {
		if err := w.send(ctx, delivery); err != nil {
			w.logger.Printf("notification delivery %d: %v", delivery.ID, err)
		}
	}
}

func (w *Worker) send(ctx context.Context, delivery Delivery) error {
	result, err := w.sender.Send(ctx, delivery)
	if err != nil {
		return w.retryOrFail(ctx, delivery, err.Error())
	}

	switch {
	case result.StatusCode >= http.StatusOK && result.StatusCode < http.StatusMultipleChoices:
		if err := w.repo.MarkSent(ctx, delivery.ID); err != nil {
			return fmt.Errorf("mark sent: %w", err)
		}
		return nil

	case result.StatusCode == http.StatusNotFound || result.StatusCode == http.StatusGone:
		if err := w.repo.DeactivateSubscription(ctx, delivery.SubscriptionID); err != nil {
			return fmt.Errorf("deactivate subscription: %w", err)
		}
		if err := w.repo.MarkFailed(ctx, delivery.ID, "push subscription expired"); err != nil {
			return fmt.Errorf("mark expired delivery: %w", err)
		}
		return nil

	case result.StatusCode == http.StatusTooManyRequests || result.StatusCode >= http.StatusInternalServerError:
		return w.retryOrFail(
			ctx,
			delivery,
			fmt.Sprintf("push provider returned %d", result.StatusCode),
		)

	default:
		message := fmt.Sprintf("push provider returned permanent status %d", result.StatusCode)
		if err := w.repo.MarkFailed(ctx, delivery.ID, message); err != nil {
			return fmt.Errorf("mark failed: %w", err)
		}
		return nil
	}
}

func (w *Worker) retryOrFail(ctx context.Context, delivery Delivery, message string) error {
	if delivery.Attempts >= maxAttempts {
		if err := w.repo.MarkFailed(ctx, delivery.ID, message); err != nil {
			return fmt.Errorf("mark exhausted delivery: %w", err)
		}
		return nil
	}

	if err := w.repo.MarkRetry(ctx, Retry{
		DeliveryID: delivery.ID,
		NextTry:    retryAt(delivery.Attempts),
		Message:    message,
	}); err != nil {
		return fmt.Errorf("schedule retry: %w", err)
	}
	return nil
}
