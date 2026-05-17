#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
if docker info >/dev/null 2>&1; then
  dc() { docker compose "$@"; }
else
  dc() { sudo -E docker compose "$@"; }
fi
dc pull api
dc up -d --no-deps --force-recreate api
dc ps
