#!/usr/bin/env bash
set -eo pipefail

cd "$(dirname "$0")/.."

if [ ! -f ".env" ]; then
  echo "Error: .env file not found in $(pwd)"
  exit 1
fi

export $(grep -v '^#' .env | xargs)

go run cmd/server/main.go \
  --emailjs-service-id="${EMAILJS_SERVICE_ID}" \
  --emailjs-template-id="${EMAILJS_TEMPLATE_ID}" \
  --emailjs-public-key="${EMAILJS_PUBLIC_KEY}" \
  --emailjs-private-key="${EMAILJS_PRIVATE_KEY}" \
  --google-measurement-id="${GOOGLE_MEASUREMENT_ID}"
