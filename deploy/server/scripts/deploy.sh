#!/usr/bin/env bash
# Deploy: faqat API image yangilanadi. DB/Redis volume'lari va ichidagi data saqlanadi.
# Migration: faqat yangi .up.sql (migrate up), mavjud qatorlar o'chirilmaydi.
set -euo pipefail

cd "$(dirname "$0")/.."

if [ ! -f .env ]; then
  echo "Xato: .env topilmadi. cp env.example .env && nano .env"
  exit 1
fi

if [ -n "${GHCR_TOKEN:-}" ] && [ -n "${GHCR_USER:-}" ]; then
  echo "$GHCR_TOKEN" | docker login ghcr.io -u "$GHCR_USER" --password-stdin
fi

# Volume (yangi nom yoki eski compose default: currency-exchange_pgdata)
for vol in currency-exchange_pgdata currency-exchange_redisdata \
           currency-exchange-pgdata currency-exchange-redisdata; do
  if docker volume inspect "$vol" >/dev/null 2>&1; then
    echo "Volume mavjud: $vol"
  fi
done

db_running() {
  local id
  id="$(docker compose ps -q db 2>/dev/null || true)"
  [ -n "$id" ] && [ "$(docker inspect -f '{{.State.Running}}' "$id" 2>/dev/null)" = "true" ]
}

if ! db_running; then
  echo "Birinchi deploy: db, redis, api ishga tushirilmoqda..."
  docker compose pull
  docker compose up -d
else
  echo "API yangilanmoqda (db/redis va datalar o'zgarmaydi)..."
  docker compose pull api
  docker compose up -d --no-deps --force-recreate api
fi

docker image prune -f
docker compose ps
