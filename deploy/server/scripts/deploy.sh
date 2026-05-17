#!/usr/bin/env bash
# Deploy: faqat API image yangilanadi. DB/Redis volume'lari va ichidagi data saqlanadi.
set -euo pipefail

cd "$(dirname "$0")/.."

if [ ! -f .env ]; then
  echo "Xato: .env topilmadi. cp env.example .env && nano .env"
  exit 1
fi

# GitHub Actions user ko'pincha docker guruhida emas — sudo fallback
if docker info >/dev/null 2>&1; then
  dk() { docker "$@"; }
  dc() { docker compose "$@"; }
else
  echo "Docker socket: sudo orqali ishlatiladi (sudo usermod -aG docker \$USER tavsiya etiladi)"
  dk() { sudo -E docker "$@"; }
  dc() { sudo -E docker compose "$@"; }
fi

if [ -n "${GHCR_TOKEN:-}" ] && [ -n "${GHCR_USER:-}" ]; then
  echo "$GHCR_TOKEN" | dk login ghcr.io -u "$GHCR_USER" --password-stdin
fi

for vol in currency-exchange_pgdata currency-exchange_redisdata \
           currency-exchange-pgdata currency-exchange-redisdata; do
  if dk volume inspect "$vol" >/dev/null 2>&1; then
    echo "Volume mavjud: $vol"
  fi
done

db_running() {
  local id
  id="$(dc ps -q db 2>/dev/null || true)"
  [ -n "$id" ] && [ "$(dk inspect -f '{{.State.Running}}' "$id" 2>/dev/null)" = "true" ]
}

if ! db_running; then
  echo "Birinchi deploy: db, redis, api ishga tushirilmoqda..."
  dc pull
  dc up -d
else
  echo "API yangilanmoqda (db/redis va datalar o'zgarmaydi)..."
  dc pull api
  dc up -d --no-deps --force-recreate api
fi

dk image prune -f
dc ps
