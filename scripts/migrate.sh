#!/bin/sh
set -eu

: "${DATABASE_URL:?DATABASE_URL is required}"

log_db_target() {
    db_host="$(printf '%s' "$DATABASE_URL" | sed -n 's|^postgres://[^@]*@\([^:/]*\).*|\1|p')"
    db_name="$(printf '%s' "$DATABASE_URL" | sed -n 's|^postgres://[^/]*/\([^?]*\).*|\1|p')"
    echo "Migration target: host=${db_host:-unknown} db=${db_name:-unknown}"
}

wait_for_db() {
    tries=0
    max=30

    while [ "$tries" -lt "$max" ]; do
        if psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -c 'SELECT 1' >/dev/null 2>&1; then
            return 0
        fi

        tries=$((tries + 1))
        echo "Waiting for database... ($tries/$max)"
        sleep 2
    done

    echo "ERROR: could not connect to database" >&2
    log_db_target >&2
    echo "Check DATABASE_URL: host must match compose alias (prod-db/dev-db or db), password must match POSTGRES_PASSWORD" >&2
    psql "$DATABASE_URL" -c 'SELECT 1' 2>&1 || true
    exit 2
}

log_db_target
wait_for_db

psql "$DATABASE_URL" -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
SQL

for migration in /migrations/*.sql; do
    [ -f "$migration" ] || continue

    version="$(basename "$migration" .sql)"
    applied="$(psql "$DATABASE_URL" -tAc "SELECT 1 FROM schema_migrations WHERE version = '$version'")"

    if [ "$applied" = "1" ]; then
        echo "Skipping migration $version (already applied)"
        continue
    fi

    echo "Applying migration $version"
    psql "$DATABASE_URL" -v ON_ERROR_STOP=1 --single-transaction \
        -f "$migration" \
        -c "INSERT INTO schema_migrations (version) VALUES ('$version')"
done

echo "Migrations complete"
