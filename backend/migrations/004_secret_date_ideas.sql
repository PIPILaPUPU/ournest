ALTER TABLE date_ideas
    ADD COLUMN IF NOT EXISTS is_secret BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_date_ideas_secret
    ON date_ideas (is_secret);
