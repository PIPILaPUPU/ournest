#!/bin/sh
set -eu

: "${DATABASE_URL:?DATABASE_URL is required}"

psql "$DATABASE_URL" -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
SQL

for migration in /migrations/*.sql; do
    version="$(basename "$migration" .sql)"
    applied="$(psql "$DATABASE_URL" -tAc "SELECT 1 FROM schema_migrations WHERE version = '$version'")"

    if [ "$applied" = "1" ]; then
        continue
    fi

    echo "Applying migration $version"
    psql "$DATABASE_URL" -v ON_ERROR_STOP=1 --single-transaction \
        -f "$migration" \
        -c "INSERT INTO schema_migrations (version) VALUES ('$version')"
done
