#!/usr/bin/env bash
set -euo pipefail

ACTION="${1:-up}"
ENV_FILE="${2:-configs/.env}"

if [ ! -f "$ENV_FILE" ]; then
  echo "env file not found: $ENV_FILE (copy configs/.env.example to configs/.env)"
  exit 1
fi

set -a
source "$ENV_FILE"
set +a

# -path is shorthand for -source=file://path, so pass the bare directory.
MIGRATION_URL="migrations"
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

# Prefer the official `migrate` CLI binary which ships with all database drivers.
# Install it once via: `brew install golang-migrate` (or download from GitHub releases).
# Fallback: `go run` with driver side-effects registered by tools/tools.go (build tag "tools").
if command -v migrate >/dev/null 2>&1; then
  migrate -path "$MIGRATION_URL" -database "$DB_URL" "$ACTION"
else
  go run -tags tools github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
    -path "$MIGRATION_URL" \
    -database "$DB_URL" \
    "$ACTION"
fi
