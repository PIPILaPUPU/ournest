CREATE TABLE IF NOT EXISTS push_subscriptions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    endpoint    TEXT        NOT NULL UNIQUE,
    p256dh      TEXT        NOT NULL,
    auth_key    TEXT        NOT NULL,
    user_agent  TEXT,
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_push_subscriptions_user_active
    ON push_subscriptions (user_id, is_active);

CREATE TABLE IF NOT EXISTS notification_outbox (
    id             BIGSERIAL PRIMARY KEY,
    event_type     VARCHAR(100) NOT NULL,
    actor_user_id  BIGINT       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    recipient_id   BIGINT       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    payload        JSONB        NOT NULL,
    status         VARCHAR(20)  NOT NULL DEFAULT 'pending',
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    dispatched_at  TIMESTAMPTZ,

    CONSTRAINT notification_outbox_status_check
        CHECK (status IN ('pending', 'dispatched'))
);

CREATE INDEX IF NOT EXISTS idx_notification_outbox_pending
    ON notification_outbox (created_at)
    WHERE status = 'pending';

CREATE TABLE IF NOT EXISTS notification_deliveries (
    id               BIGSERIAL PRIMARY KEY,
    outbox_id        BIGINT       NOT NULL REFERENCES notification_outbox (id) ON DELETE CASCADE,
    subscription_id  BIGINT       NOT NULL REFERENCES push_subscriptions (id) ON DELETE CASCADE,
    status           VARCHAR(20)  NOT NULL DEFAULT 'pending',
    attempts         INTEGER      NOT NULL DEFAULT 0,
    next_attempt_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    last_attempt_at  TIMESTAMPTZ,
    last_error       TEXT,
    sent_at          TIMESTAMPTZ,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT notification_delivery_unique UNIQUE (outbox_id, subscription_id),
    CONSTRAINT notification_delivery_status_check
        CHECK (status IN ('pending', 'processing', 'retry', 'sent', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_notification_deliveries_ready
    ON notification_deliveries (next_attempt_at)
    WHERE status IN ('pending', 'retry', 'processing');
