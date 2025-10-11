#!/usr/bin/env bash
set -eo pipefail
cd "$(dirname "$0")/.." || exit 1

if [ ! -f ".env" ]; then
  echo "Error: .env file not found in $(pwd)"
  exit 1
fi

set -a
source .env
set +a

CMD=${1:-up}
COUNT=${2:-}

if [ -z "$DATABASE_URL" ]; then
  echo "Error: DATABASE_URL not set in .env"
  exit 1
fi

echo "Running migrations: $CMD $COUNT"
migrate -path internal/database/migrations -database "$DATABASE_URL" "$CMD" $COUNT
