-- Схема для приложения wishlist (два пользователя-партнёра)

DO $$
BEGIN
    CREATE TYPE wish_status AS ENUM ('wanted', 'reserved', 'done');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wishes (
    id           BIGSERIAL PRIMARY KEY,
    owner_id     BIGINT       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title        VARCHAR(255) NOT NULL,
    description  TEXT,
    url          TEXT,
    price        NUMERIC(10, 2),
    status       wish_status  NOT NULL DEFAULT 'wanted',
    group_name   VARCHAR(100) NOT NULL DEFAULT 'Общее',
    group_color  VARCHAR(20)  NOT NULL DEFAULT 'slate',
    reserved_by  BIGINT       REFERENCES users (id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT wishes_reserved_requires_user
        CHECK (
            (status = 'reserved' AND reserved_by IS NOT NULL)
            OR (status != 'reserved')
        )
);

CREATE INDEX IF NOT EXISTS idx_wishes_owner_id ON wishes (owner_id);
CREATE INDEX IF NOT EXISTS idx_wishes_status ON wishes (status);

-- Тестовые пользователи (пароль для обоих: secret)
-- bcrypt hash of "secret"
INSERT INTO users (username, password_hash) VALUES
    ('alice', '$2a$10$mUO4ptNAp9vsvaULWt.rpOibC6RepIqjmg8z1PAHSOQ1P6DWkmi1W'),
    ('bob',   '$2a$10$mUO4ptNAp9vsvaULWt.rpOibC6RepIqjmg8z1PAHSOQ1P6DWkmi1W')
ON CONFLICT (username) DO NOTHING;

INSERT INTO wishes (owner_id, title, description, url, price, group_name, group_color)
SELECT seed.*
FROM (
    VALUES
        (1::BIGINT, 'Книга по Go', 'Для освежения знаний', 'https://example.com/go-book', 45.00::NUMERIC, 'Учеба', 'violet'),
        (2::BIGINT, 'Наушники', NULL, 'https://example.com/headphones', 120.00::NUMERIC, 'Техника', 'blue')
) AS seed(owner_id, title, description, url, price, group_name, group_color)
WHERE NOT EXISTS (SELECT 1 FROM wishes);
