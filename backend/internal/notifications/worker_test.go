package notifications

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"testing"
)

type fakeWorkerRepository struct {
	sent                  bool
	retried               bool
	failed                bool
	deactivated           bool
	deactivatedSubID      int64
	lastRetry             Retry
	lastFailureDeliveryID int64
}

func (r *fakeWorkerRepository) DispatchBatch(context.Context, int) (int, error) {
	return 0, nil
}

func (r *fakeWorkerRepository) ClaimDeliveries(context.Context, int) ([]Delivery, error) {
	return nil, nil
}

func (r *fakeWorkerRepository) MarkSent(context.Context, int64) error {
	r.sent = true
	return nil
}

func (r *fakeWorkerRepository) MarkRetry(_ context.Context, retry Retry) error {
	r.retried = true
	r.lastRetry = retry
	return nil
}

func (r *fakeWorkerRepository) MarkFailed(_ context.Context, id int64, _ string) error {
	r.failed = true
	r.lastFailureDeliveryID = id
	return nil
}

func (r *fakeWorkerRepository) DeactivateSubscription(_ context.Context, id int64) error {
	r.deactivated = true
	r.deactivatedSubID = id
	return nil
}

type fakeSender struct {
	result SendResult
	err    error
}

func (s fakeSender) Send(context.Context, Delivery) (SendResult, error) {
	return s.result, s.err
}

func TestWorkerDeactivatesExpiredSubscription(t *testing.T) {
	repo := &fakeWorkerRepository{}
	worker := NewWorker(
		repo,
		fakeSender{result: SendResult{StatusCode: http.StatusGone}},
		log.New(io.Discard, "", 0),
	)

	err := worker.send(context.Background(), Delivery{ID: 11, SubscriptionID: 22})
	if err != nil {
		t.Fatalf("send returned error: %v", err)
	}
	if !repo.deactivated || repo.deactivatedSubID != 22 {
		t.Fatal("expired subscription was not deactivated")
	}
	if !repo.failed || repo.lastFailureDeliveryID != 11 {
		t.Fatal("expired delivery was not marked failed")
	}
	if repo.retried {
		t.Fatal("expired subscription must not be retried")
	}
}

func TestWorkerRetriesTransientFailure(t *testing.T) {
	repo := &fakeWorkerRepository{}
	worker := NewWorker(
		repo,
		fakeSender{err: errors.New("network unavailable")},
		log.New(io.Discard, "", 0),
	)

	err := worker.send(context.Background(), Delivery{ID: 33, Attempts: 1})
	if err != nil {
		t.Fatalf("send returned error: %v", err)
	}
	if !repo.retried || repo.lastRetry.DeliveryID != 33 {
		t.Fatal("transient failure was not scheduled for retry")
	}
	if repo.failed {
		t.Fatal("delivery failed before attempts were exhausted")
	}
}

func TestWorkerStopsAfterMaximumAttempts(t *testing.T) {
	repo := &fakeWorkerRepository{}
	worker := NewWorker(
		repo,
		fakeSender{result: SendResult{StatusCode: http.StatusInternalServerError}},
		log.New(io.Discard, "", 0),
	)

	err := worker.send(context.Background(), Delivery{ID: 44, Attempts: maxAttempts})
	if err != nil {
		t.Fatalf("send returned error: %v", err)
	}
	if !repo.failed || repo.retried {
		t.Fatal("exhausted delivery must be failed without another retry")
	}
}
