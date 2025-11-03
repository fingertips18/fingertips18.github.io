#!/usr/bin/env bash
set -eo pipefail
cd "$(dirname "$0")/.."

go build -o cmd/tmp/main ./cmd/server
