#!/usr/bin/env bash
set -euo pipefail

ENVIRONMENT="${1:?usage: ./scripts/deploy.sh prod|dev}"
REPO_DIR="${REPO_DIR:-$HOME/ournest}"

cd "$REPO_DIR"

ensure_network() {
  docker network inspect ournest_shared >/dev/null 2>&1 || docker network create ournest_shared
}

remove_container_if_exists() {
  local name="$1"
  if docker container inspect "$name" >/dev/null 2>&1; then
    echo "Removing stale container: $name"
    docker rm -f "$name"
  fi
}

cleanup_prod_stale_containers() {
  remove_container_if_exists wishlist_caddy
  remove_container_if_exists wishlist_backend
  remove_container_if_exists wishlist_frontend
  remove_container_if_exists wishlist_db
}

compose_down_legacy_projects() {
  local env_file="$1"
  local compose_file="$2"

  docker compose -p ournest --env-file "$env_file" -f "$compose_file" down --remove-orphans 2>/dev/null || true
  docker compose --env-file "$env_file" -f "$compose_file" down --remove-orphans 2>/dev/null || true
}

compose_up() {
  local env_file="$1"
  local compose_file="$2"

  compose_down_legacy_projects "$env_file" "$compose_file"
  docker compose --env-file "$env_file" -f "$compose_file" up -d --build --remove-orphans
}

case "$ENVIRONMENT" in
  prod)
    git fetch origin
    git checkout main
    git pull origin main
    ensure_network
    cleanup_prod_stale_containers
    compose_up .env.production docker-compose.prod.yml
    ;;
  dev)
    git fetch origin
    git checkout dev
    git pull origin dev
    ensure_network
    compose_up .env.development docker-compose.dev.yml
    ;;
  *)
    echo "unknown environment: $ENVIRONMENT (expected prod or dev)" >&2
    exit 1
    ;;
esac

docker image prune -f
