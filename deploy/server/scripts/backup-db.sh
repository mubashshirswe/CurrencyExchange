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

# shellcheck disable=SC1091
source "$(dirname "$0")/lib/load-db-env.sh"
load_db_env

STAMP="$(date +%Y%m%d-%H%M%S)"
OUT="backups/db-${STAMP}.sql.gz"

echo "Backup: ${POSTGRES_DB} -> ${OUT}"
docker compose exec -T db pg_dump -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" --no-owner --no-acl \
  | gzip -9 > "${OUT}"

ls -lh "${OUT}"
echo "Tayyor. Yangi serverga scp qiling."
