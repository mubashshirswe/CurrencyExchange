#!/bin/sh
set -e

if [ "${RUN_MIGRATIONS:-true}" = "true" ] && [ -n "${DB_ADDR:-}" ]; then
  echo "Running database migrations..."
  migrate -path=/migrations -database="${DB_ADDR}" up
fi

exec "$@"
