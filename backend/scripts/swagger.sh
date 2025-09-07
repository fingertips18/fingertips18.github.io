#!/usr/bin/env bash
set -eo pipefail

# Move to project root (assuming script is in scripts/ or similar)
cd "$(dirname "$0")/.."

# Generate Swagger docs
swag init -g cmd/server/main.go