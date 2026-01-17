#!/usr/bin/env bash
set -eo pipefail
cd "$(dirname "$0")/.." || exit 1

MIGRATION_NAME=$1

if [ -z "$MIGRATION_NAME" ]; then
  echo "Error: Migration name required"
  echo "Usage: ./scripts/create-migration.sh <migration_name>"
  echo "Example: ./scripts/create-migration.sh fix_file_table_schema"
  exit 1
fi

TIMESTAMP=$(date -u +%Y%m%d%H%M%S)
MIGRATION_DIR="internal/database/migrations"
UP_FILE="${MIGRATION_DIR}/${TIMESTAMP}_${MIGRATION_NAME}.up.sql"
DOWN_FILE="${MIGRATION_DIR}/${TIMESTAMP}_${MIGRATION_NAME}.down.sql"

# Create the up migration file
cat > "$UP_FILE" << EOF
-- Migration: ${MIGRATION_NAME}
-- Created: $(date -u)

EOF

# Create the down migration file
cat > "$DOWN_FILE" << EOF
-- Rollback: ${MIGRATION_NAME}
-- Created: $(date -u)

EOF

echo "✅ Created migration files:"
echo "   $UP_FILE"
echo "   $DOWN_FILE"