#!/usr/bin/env bash
# Yangi serverda: ./scripts/restore-db.sh backups/db-XXXX.sql.gz
set -euo pipefail

cd "$(dirname "$0")/.."

DUMP="${1:-}"
if [ -z "${DUMP}" ] || [ ! -f "${DUMP}" ]; then
  echo "Foydalanish: $0 backups/db-YYYYMMDD-HHMMSS.sql.gz"
  exit 1
fi

if [ ! -f .env ]; then
  echo "Xato: .env yo'q (yangi serverning o'z .env — parol/user eskisi bilan bir xil bo'lishi shart emas)"
  exit 1
fi

set -a
# shellcheck disable=SC1091
source .env
set +a

# shellcheck disable=SC1091
source "$(dirname "$0")/lib/load-db-env.sh"
load_db_env

echo "DB va Redis ishga tushirilmoqda (API hali yo'q)..."
docker compose up -d db redis

echo "DB healthy kutilyapti..."
for i in $(seq 1 30); do
  if docker compose exec -T db pg_isready -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" >/dev/null 2>&1; then
    break
  fi
  sleep 2
done

echo "Restore: ${DUMP}"
gunzip -c "${DUMP}" | docker compose exec -T db psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -v ON_ERROR_STOP=1

echo "Restore tugadi. Endi: docker compose up -d"
