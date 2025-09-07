#!/usr/bin/env bash
set -eo pipefail
cd "$(dirname "$0")/.."

if [ ! -f ".env" ]; then
  echo "Error: .env file not found in $(pwd)"
  exit 1
fi

set -a
source .env
set +a

migrate -database $POSTGRES_CONNECTION_STRING -path internal/database/migrations up