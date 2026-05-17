#!/usr/bin/env bash
# Birinchi marta serverda: bash scripts/bootstrap.sh
set -euo pipefail

cd "$(dirname "$0")/.."

if ! command -v docker >/dev/null 2>&1; then
  echo "Docker topilmadi. O'rnating: https://docs.docker.com/engine/install/"
  exit 1
fi

if [ ! -f .env ]; then
  cp env.example .env
  echo ".env yaratildi — parollarni tahrirlang: nano .env"
  echo "Keyin: ./scripts/deploy.sh"
  exit 0
fi

echo "Ogoh: docker compose down -v volume'larni o'chiradi — userlar va DB data yo'qoladi!"
chmod +x scripts/deploy.sh
./scripts/deploy.sh
