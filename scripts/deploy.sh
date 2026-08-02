#!/usr/bin/env bash
set -euo pipefail

ENVIRONMENT="${1:?usage: ./scripts/deploy.sh prod|dev}"
REPO_DIR="${REPO_DIR:-$HOME/ournest}"

cd "$REPO_DIR"

ensure_network() {
  docker network inspect ournest_shared >/dev/null 2>&1 || docker network create ournest_shared
}

case "$ENVIRONMENT" in
  prod)
    git fetch origin
    git checkout main
    git pull origin main
    ensure_network
    docker compose --env-file .env.production -f docker-compose.prod.yml up -d --build
    ;;
  dev)
    git fetch origin
    git checkout dev
    git pull origin dev
    ensure_network
    docker compose --env-file .env.development -f docker-compose.dev.yml up -d --build
    ;;
  *)
    echo "unknown environment: $ENVIRONMENT (expected prod or dev)" >&2
    exit 1
    ;;
esac

docker image prune -f
