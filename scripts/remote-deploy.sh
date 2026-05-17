#!/usr/bin/env bash
# Eski yo'l — kanonik skript: deploy/server/scripts/deploy.sh
set -euo pipefail
DIR="$(cd "$(dirname "$0")/.." && pwd)"
exec "$DIR/deploy/server/scripts/deploy.sh"
