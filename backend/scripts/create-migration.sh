#!/bin/bash
set -eo pipefail

# Change to repository root
cd "$(git rev-parse --show-toplevel)/backend" || exit 1

# Validate migration name is provided
if [ -z "$1" ]; then
  echo "Error: Migration name is required"
  echo "Usage: ./scripts/create-migration.sh <migration_name>"
  echo "Example: ./scripts/create-migration.sh add_user_table"
  exit 1
fi

MIGRATION_NAME="$1"
MIGRATIONS_DIR="internal/database/migrations"

# Ensure migrations directory exists
mkdir -p "$MIGRATIONS_DIR"

# Create the migration files
echo "Creating migration: $MIGRATION_NAME"
migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$MIGRATION_NAME"

echo "✓ Migration files created successfully"
echo "  Edit the generated files in $MIGRATIONS_DIR/"