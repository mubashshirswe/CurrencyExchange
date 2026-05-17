#!/usr/bin/env bash
# Eski serverda: ./scripts/backup-db.sh
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

STAMP="$(date +%Y%m%d-%H%M%S)"
OUT="backups/db-${STAMP}.sql.gz"

echo "Backup: ${POSTGRES_DB} -> ${OUT}"
docker compose exec -T db pg_dump -U "${POSTGRES_USER:-app}" -d "${POSTGRES_DB:-currency_exchange}" --no-owner --no-acl \
  | gzip -9 > "${OUT}"

ls -lh "${OUT}"
echo "Tayyor. Yangi serverga scp qiling."
