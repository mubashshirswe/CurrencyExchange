#!/usr/bin/env bash
# Yangi serverda DB restore (bitta fayl, lib kerak emas).
#
# Oldin:
#   1) Dump ni backups/ ga qo'ying:
#        mkdir -p backups
#        # Mac: scp root@ESKI:.../backups/db-XXX.sql.gz user@YANGI:/app/backend/currency-exchange/backups/
#   2) .env tayyor: cp env.example .env && nano .env
#        (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, API_IMAGE, JWTSECRET)
#
# Ishga tushirish:
#   chmod +x restore.sh
#   ./restore.sh --fresh backups/db-20260518-032449.sql.gz
#   (--fresh: mavjud jadvallarni o'chirib, dumpdan to'liq tiklaydi)
#
# Keyin:
#   docker compose up -d
set -euo pipefail

cd "$(dirname "$0")"

FRESH=0
DUMP=""
while [ $# -gt 0 ]; do
  case "$1" in
    --fresh) FRESH=1; shift ;;
    -h|--help)
      echo "Foydalanish: $0 [--fresh] backups/db-YYYYMMDD-HHMMSS.sql.gz"
      exit 0
      ;;
    *)
      DUMP="$1"
      shift
      ;;
  esac
done

if [ -z "${DUMP}" ] || [ ! -f "${DUMP}" ]; then
  echo "Foydalanish: $0 [--fresh] backups/db-YYYYMMDD-HHMMSS.sql.gz"
  exit 1
fi

if [ ! -f .env ]; then
  echo "Xato: .env yo'q. cp env.example .env && nano .env"
  exit 1
fi

set -a
# shellcheck disable=SC1091
source .env
set +a

# POSTGRES_* yoki DB_ADDR dan o'qish
if [ -z "${POSTGRES_DB:-}" ] || [ -z "${POSTGRES_USER:-}" ]; then
  if [ -z "${DB_ADDR:-}" ]; then
    echo "Xato: .env da POSTGRES_USER/POSTGRES_DB yoki DB_ADDR kerak"
    exit 1
  fi
  if ! command -v python3 >/dev/null 2>&1; then
    echo "Xato: DB_ADDR uchun python3 kerak"
    exit 1
  fi
  eval "$(DB_ADDR="$DB_ADDR" python3 - <<'PY'
import os, shlex
from urllib.parse import urlparse, unquote
u = urlparse(os.environ["DB_ADDR"])
user = unquote(u.username or "") or "app"
password = unquote(u.password or "")
db = (u.path or "/").lstrip("/").split("?")[0] or "currency_exchange"
for k, v in [("POSTGRES_USER", user), ("POSTGRES_PASSWORD", password), ("POSTGRES_DB", db)]:
    print(f"export {k}={shlex.quote(v)}")
PY
)"
fi

POSTGRES_USER="${POSTGRES_USER:-app}"
POSTGRES_DB="${POSTGRES_DB:-currency_exchange}"

find_db_container() {
  local cid project="${COMPOSE_PROJECT:-currency-exchange}"
  cid="$(docker ps -q \
    --filter "label=com.docker.compose.project=${project}" \
    --filter "label=com.docker.compose.service=db" 2>/dev/null | head -1)"
  [ -n "$cid" ] && { echo "$cid"; return 0; }
  cid="$(docker ps -q --filter "ancestor=postgres:16-alpine" 2>/dev/null | head -1)"
  [ -n "$cid" ] && { echo "$cid"; return 0; }
  echo "Xato: Postgres konteyner topilmadi. Avval: docker compose up -d db" >&2
  return 1
}

db_exec() {
  docker exec -i "$(find_db_container)" "$@"
}

if docker info >/dev/null 2>&1; then
  dc() { docker compose "$@"; }
else
  dc() { sudo docker compose "$@"; }
fi

echo ">> DB va Redis ishga tushirilmoqda..."
dc up -d db redis

echo ">> DB tayyor bo'lishi kutilmoqda..."
for _ in $(seq 1 30); do
  if db_exec pg_isready -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" >/dev/null 2>&1; then
    break
  fi
  sleep 2
done

if [ "${FRESH}" -eq 1 ]; then
  echo ">> --fresh: mavjud jadvallar o'chirilmoqda..."
  db_exec psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -v ON_ERROR_STOP=1 <<-SQL
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
		GRANT ALL ON SCHEMA public TO ${POSTGRES_USER};
		GRANT ALL ON SCHEMA public TO public;
	SQL
fi

echo ">> Restore: ${DUMP} -> ${POSTGRES_DB} (user: ${POSTGRES_USER})"
gunzip -c "${DUMP}" | db_exec psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -v ON_ERROR_STOP=1

echo ""
echo ">> Restore tugadi."
echo ">> Endi: docker compose up -d"
