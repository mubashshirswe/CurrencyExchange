#!/usr/bin/env bash
# .env yuklangandan keyin source qiling.
# POSTGRES_* yo'q, DB_ADDR bor bo'lsa — user/password/db nomini ajratadi.

load_db_env() {
  if [ -n "${POSTGRES_DB:-}" ] && [ -n "${POSTGRES_USER:-}" ]; then
    POSTGRES_USER="${POSTGRES_USER:-app}"
    POSTGRES_DB="${POSTGRES_DB:-currency_exchange}"
    return 0
  fi

  if [ -z "${DB_ADDR:-}" ]; then
    echo "Xato: .env da POSTGRES_DB yoki DB_ADDR kerak"
    return 1
  fi

  if ! command -v python3 >/dev/null 2>&1; then
    echo "Xato: DB_ADDR dan o'qish uchun python3 kerak"
    return 1
  fi

  eval "$(DB_ADDR="$DB_ADDR" python3 - <<'PY'
import os, shlex
from urllib.parse import urlparse, unquote

addr = os.environ.get("DB_ADDR", "")
u = urlparse(addr)
if u.scheme not in ("postgres", "postgresql"):
    raise SystemExit("DB_ADDR postgres URL emas")

user = unquote(u.username or "") or "app"
password = unquote(u.password or "")
db = (u.path or "/").lstrip("/").split("?")[0] or "currency_exchange"

for key, val in (
    ("POSTGRES_USER", user),
    ("POSTGRES_PASSWORD", password),
    ("POSTGRES_DB", db),
):
    print(f"export {key}={shlex.quote(val)}")
PY
)"
}
