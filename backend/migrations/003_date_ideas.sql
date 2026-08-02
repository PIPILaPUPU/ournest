DO $$
BEGIN
    CREATE TYPE date_idea_status AS ENUM ('planned', 'done');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS date_ideas (
    id          BIGSERIAL PRIMARY KEY,
    author_id   BIGINT           NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title       VARCHAR(255)     NOT NULL,
    description TEXT,
    status      date_idea_status NOT NULL DEFAULT 'planned',
    created_at  TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_date_ideas_author_id ON date_ideas (author_id);
CREATE INDEX IF NOT EXISTS idx_date_ideas_status ON date_ideas (status);
