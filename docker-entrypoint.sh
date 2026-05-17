#!/bin/sh
set -e

if [ "${RUN_MIGRATIONS:-true}" = "true" ]; then
  DB_URL="$(./currency_exchange -print-db-url)"
  echo "Running database migrations..."
  migrate -path=/migrations -database="${DB_URL}" up
fi

exec "$@"
