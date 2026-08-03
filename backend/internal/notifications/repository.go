package notifications

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db      *pgxpool.Pool
	enabled bool
}

func NewRepository(db *pgxpool.Pool, enabled bool) *Repository {
	return &Repository{db: db, enabled: enabled}
}

func (r *Repository) SaveSubscription(
	ctx context.Context,
	userID int,
	request SubscribeRequest,
	userAgent string,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO push_subscriptions (
			user_id, endpoint, p256dh, auth_key, user_agent
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (endpoint) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			p256dh = EXCLUDED.p256dh,
			auth_key = EXCLUDED.auth_key,
			user_agent = EXCLUDED.user_agent,
			is_active = TRUE,
			updated_at = NOW()`,
		userID,
		request.Endpoint,
		request.Keys.P256DH,
		request.Keys.Auth,
		userAgent,
	)
	if err != nil {
		return fmt.Errorf("save push subscription: %w", err)
	}
	return nil
}

func (r *Repository) RemoveSubscription(ctx context.Context, userID int, endpoint string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE push_subscriptions
		SET is_active = FALSE, updated_at = NOW()
		WHERE user_id = $1 AND endpoint = $2`,
		userID,
		endpoint,
	)
	if err != nil {
		return fmt.Errorf("remove push subscription: %w", err)
	}
	return nil
}

func (r *Repository) EnqueueForPartner(
	ctx context.Context,
	tx pgx.Tx,
	eventType string,
	actorUserID int,
	payload EventPayload,
) error {
	if !r.enabled {
		return nil
	}

	var recipientID int
	err := tx.QueryRow(ctx, `
		SELECT id
		FROM users
		WHERE id <> $1
		ORDER BY id
		LIMIT 1`,
		actorUserID,
	).Scan(&recipientID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("find notification recipient: %w", err)
	}

	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal notification payload: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO notification_outbox (
			event_type, actor_user_id, recipient_id, payload
		)
		VALUES ($1, $2, $3, $4)`,
		eventType,
		actorUserID,
		recipientID,
		rawPayload,
	)
	if err != nil {
		return fmt.Errorf("enqueue notification: %w", err)
	}
	return nil
}

func (r *Repository) DispatchBatch(ctx context.Context, limit int) (int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin notification dispatch: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	rows, err := tx.Query(ctx, `
		SELECT id, recipient_id
		FROM notification_outbox
		WHERE status = 'pending'
		ORDER BY created_at
		FOR UPDATE SKIP LOCKED
		LIMIT $1`,
		limit,
	)
	if err != nil {
		return 0, fmt.Errorf("select notification outbox: %w", err)
	}

	type event struct {
		id          int64
		recipientID int
	}
	events := make([]event, 0, limit)
	for rows.Next() {
		var item event
		if err := rows.Scan(&item.id, &item.recipientID); err != nil {
			rows.Close()
			return 0, fmt.Errorf("scan notification outbox: %w", err)
		}
		events = append(events, item)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return 0, fmt.Errorf("iterate notification outbox: %w", err)
	}
	rows.Close()

	for _, item := range events {
		if _, err := tx.Exec(ctx, `
			INSERT INTO notification_deliveries (outbox_id, subscription_id)
			SELECT $1, id
			FROM push_subscriptions
			WHERE user_id = $2 AND is_active = TRUE
			ON CONFLICT (outbox_id, subscription_id) DO NOTHING`,
			item.id,
			item.recipientID,
		); err != nil {
			return 0, fmt.Errorf("create notification deliveries: %w", err)
		}

		if _, err := tx.Exec(ctx, `
			UPDATE notification_outbox
			SET status = 'dispatched', dispatched_at = NOW()
			WHERE id = $1`,
			item.id,
		); err != nil {
			return 0, fmt.Errorf("mark notification dispatched: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit notification dispatch: %w", err)
	}
	return len(events), nil
}

func (r *Repository) ClaimDeliveries(ctx context.Context, limit int) ([]Delivery, error) {
	rows, err := r.db.Query(ctx, `
		WITH claimed AS (
			SELECT id
			FROM notification_deliveries
			WHERE (
				status IN ('pending', 'retry')
				AND next_attempt_at <= NOW()
			) OR (
				status = 'processing'
				AND last_attempt_at < NOW() - INTERVAL '5 minutes'
			)
			ORDER BY next_attempt_at
			FOR UPDATE SKIP LOCKED
			LIMIT $1
		),
		updated AS (
			UPDATE notification_deliveries d
			SET
				status = 'processing',
				attempts = d.attempts + 1,
				last_attempt_at = NOW()
			FROM claimed c
			WHERE d.id = c.id
			RETURNING d.id, d.subscription_id, d.outbox_id, d.attempts
		)
		SELECT
			u.id,
			u.subscription_id,
			s.endpoint,
			s.p256dh,
			s.auth_key,
			o.payload,
			u.attempts
		FROM updated u
		JOIN push_subscriptions s ON s.id = u.subscription_id
		JOIN notification_outbox o ON o.id = u.outbox_id`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("claim notification deliveries: %w", err)
	}
	defer rows.Close()

	deliveries := make([]Delivery, 0, limit)
	for rows.Next() {
		var delivery Delivery
		if err := rows.Scan(
			&delivery.ID,
			&delivery.SubscriptionID,
			&delivery.Endpoint,
			&delivery.P256DH,
			&delivery.Auth,
			&delivery.Payload,
			&delivery.Attempts,
		); err != nil {
			return nil, fmt.Errorf("scan notification delivery: %w", err)
		}
		deliveries = append(deliveries, delivery)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate notification deliveries: %w", err)
	}
	return deliveries, nil
}

func (r *Repository) MarkSent(ctx context.Context, deliveryID int64) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notification_deliveries
		SET status = 'sent', sent_at = NOW(), last_error = NULL
		WHERE id = $1`,
		deliveryID,
	)
	return err
}

func (r *Repository) MarkRetry(ctx context.Context, retry Retry) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notification_deliveries
		SET status = 'retry', next_attempt_at = $2, last_error = $3
		WHERE id = $1`,
		retry.DeliveryID,
		retry.NextTry,
		retry.Message,
	)
	return err
}

func (r *Repository) MarkFailed(ctx context.Context, deliveryID int64, message string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notification_deliveries
		SET status = 'failed', last_error = $2
		WHERE id = $1`,
		deliveryID,
		message,
	)
	return err
}

func (r *Repository) DeactivateSubscription(ctx context.Context, subscriptionID int64) error {
	_, err := r.db.Exec(ctx, `
		UPDATE push_subscriptions
		SET is_active = FALSE, updated_at = NOW()
		WHERE id = $1`,
		subscriptionID,
	)
	return err
}

func retryAt(attempt int) time.Time {
	delays := []time.Duration{
		30 * time.Second,
		2 * time.Minute,
		10 * time.Minute,
		1 * time.Hour,
	}
	index := attempt - 1
	if index < 0 {
		index = 0
	}
	if index >= len(delays) {
		index = len(delays) - 1
	}
	return time.Now().Add(delays[index])
}
