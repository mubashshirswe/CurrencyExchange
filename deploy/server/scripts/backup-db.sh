#!/usr/bin/env bash
# Eski serverda: ./scripts/backup-db.sh
# Faqat DB_ADDR yoki POSTGRES_* kerak — API_IMAGE/POSTGRES_PASSWORD compose uchun shart emas.
# Natija: ../backups/db-YYYYMMDD-HHMMSS.sql.gz
set -euo pipefail

cd "$(dirname "$0")/.."
mkdir -p backups

if [ ! -f .env ]; then
  echo "Xato: .env yo'q"
  exit 1
fi

set -a
# shellcheck disable=SC1091
source .env
set +a

# shellcheck disable=SC1091
source "$(dirname "$0")/lib/load-db-env.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/lib/docker-db.sh"
load_db_env

STAMP="$(date +%Y%m%d-%H%M%S)"
OUT="backups/db-${STAMP}.sql.gz"
CID="$(find_db_container)"

echo "Backup: ${POSTGRES_DB} (konteyner ${CID:0:12}) -> ${OUT}"
db_exec pg_dump -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" --no-owner --no-acl \
  | gzip -9 > "${OUT}"

SIZE="$(wc -c < "${OUT}" | tr -d ' ')"
if [ "${SIZE}" -lt 1024 ]; then
  echo "Xato: backup juda kichik (${SIZE} bayt) — ehtimol xato. Faylni o'chiring va qayta urinib ko'ring."
  exit 1
fi

ls -lh "${OUT}"
echo "Tayyor. Yangi serverga scp qiling."
