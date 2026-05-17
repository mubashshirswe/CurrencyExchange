#!/usr/bin/env bash
# Ishlayotgan db konteyneri — docker compose butun .env ni talab qilmaydi.

COMPOSE_PROJECT="${COMPOSE_PROJECT:-currency-exchange}"

find_db_container() {
  local cid

  cid="$(docker ps -q \
    --filter "label=com.docker.compose.project=${COMPOSE_PROJECT}" \
    --filter "label=com.docker.compose.service=db" 2>/dev/null | head -1)"
  if [ -n "$cid" ]; then
    echo "$cid"
    return 0
  fi

  cid="$(docker ps -q --filter "ancestor=postgres:16-alpine" 2>/dev/null | head -1)"
  if [ -n "$cid" ]; then
    echo "$cid"
    return 0
  fi

  local id names
  while read -r id names; do
    case "$names" in
      *-db-1 | *_db_1 | *-db-*) echo "$id"; return 0 ;;
    esac
  done < <(docker ps --format '{{.ID}} {{.Names}}')

  echo "Xato: Postgres konteyner topilmadi. docker ps | grep -i db" >&2
  return 1
}

db_exec() {
  local cid
  cid="$(find_db_container)" || return 1
  docker exec -i "$cid" "$@"
}
