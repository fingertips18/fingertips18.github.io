# Database Migrations Guide

This guide covers the database migration workflow using the `golang-migrate` tool and helper scripts.

## Prerequisites

- `golang-migrate` CLI installed ([installation guide](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate))
- `.env` file in `backend/` with `DATABASE_URL` set
- PostgreSQL database running and accessible

## Scripts Overview

### 1. create-migration.sh (Creating Migrations)

**Location:** `backend/scripts/create-migration.sh`

**Purpose:** Generate new migration files with proper naming and structure.

**Usage:**
```bash
./scripts/create-migration.sh <migration_name>
```

**Examples:**
```bash
# Create a new table migration
./scripts/create-migration.sh add_user_table

# Add a column to existing table
./scripts/create-migration.sh add_email_to_users

# Modify constraints
./scripts/create-migration.sh update_project_type_constraint
```

**Output:**
- Creates two files in `internal/database/migrations/`:
  - `<timestamp>_<migration_name>.up.sql` (apply changes)
  - `<timestamp>_<migration_name>.down.sql` (revert changes)

**Notes:**
- Files are generated empty—you must add SQL content manually
- Timestamp is auto-generated for ordering and uniqueness
- Use descriptive names (snake_case recommended)

---

### 2. migrate.sh (Running Migrations)

**Location:** `backend/scripts/migrate.sh`

**Purpose:** Apply or revert database migrations.

**Usage:**
```bash
./scripts/migrate.sh [COMMAND] [COUNT]
```

**Parameters:**
- `COMMAND`: Migration command (default: `up`)
- `COUNT`: Optional number of migrations to apply/revert

---

## Common Migration Commands

### Apply All Pending Migrations
```bash
./scripts/migrate.sh up
```

### Apply Specific Number of Migrations
```bash
# Apply next 2 migrations
./scripts/migrate.sh up 2
```

### Rollback All Migrations
```bash
./scripts/migrate.sh down
```

### Rollback Specific Number of Migrations
```bash
# Rollback last 1 migration
./scripts/migrate.sh down 1
```

### Check Migration Status
```bash
./scripts/migrate.sh version
```

### Force Specific Version (Use with Caution)
```bash
./scripts/migrate.sh force 20250831062237
```

### Drop Everything (Destructive)
```bash
./scripts/migrate.sh drop
```
**Warning:** This will drop all tables and data!

---

## Migration Workflow

### Step 1: Create Migration Files
```bash
cd backend
./scripts/create-migration.sh add_education_table
```

### Step 2: Write SQL in Generated Files

**Example `.up.sql`:**
```sql
CREATE TABLE education (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    level TEXT CHECK (level IN ('elementary', 'junior-high-school', 'senior-high-school', 'college')) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
```

**Example `.down.sql`:**
```sql
DROP TABLE IF EXISTS education;
```

### Step 3: Run Migration
```bash
./scripts/migrate.sh up
```

### Step 4: Verify
```bash
# Check current version
./scripts/migrate.sh version

# Or connect to database and verify schema
psql $DATABASE_URL -c "\dt"
```

---

## Best Practices

### DO

- **Write reversible migrations:** Always implement proper `.down.sql` logic
- **Use idempotent operations:** Include `IF EXISTS` / `IF NOT EXISTS` where appropriate
- **Test both directions:** Run `up` then `down` to ensure clean reversal
- **Never mutate historical migrations:** Once applied to production, create new migrations instead
- **Use descriptive names:** `add_user_email_column` not `update_users`
- **One logical change per migration:** Keep migrations focused and atomic
- **Add comments:** Document complex transformations or business logic

### DON'T

- **Don't modify applied migrations:** Historical files should never change
- **Don't mix DDL and DML:** Separate schema changes from data migrations
- **Don't skip version control:** Always commit both `.up.sql` and `.down.sql`
- **Don't forget CHECK constraints:** Use database-level validation where possible
- **Don't use raw SQL in Go code:** Keep schema changes in migration files

---

## Example Migrations

### Create Table with Constraints
```sql
-- up
CREATE TABLE project (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    stack TEXT[] NOT NULL DEFAULT '{}',
    type TEXT CHECK (type IN ('web', 'mobile', 'game')),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- down
DROP TABLE IF EXISTS project;
```

### Add Column with Default
```sql
-- up
ALTER TABLE project 
ADD COLUMN status TEXT DEFAULT 'draft' NOT NULL;

-- down
ALTER TABLE project 
DROP COLUMN IF EXISTS status;
```

### Modify Existing Column
```sql
-- up
ALTER TABLE project 
ALTER COLUMN stack SET NOT NULL,
ALTER COLUMN stack SET DEFAULT '{}';

-- down
ALTER TABLE project 
ALTER COLUMN stack DROP NOT NULL,
ALTER COLUMN stack DROP DEFAULT;
```

### Add Index
```sql
-- up
CREATE INDEX idx_project_type ON project(type);
CREATE INDEX idx_project_created_at ON project(created_at);

-- down
DROP INDEX IF EXISTS idx_project_created_at;
DROP INDEX IF EXISTS idx_project_type;
```

---

## Troubleshooting

### "Dirty database version" Error
```bash
# Force to last known good version
./scripts/migrate.sh force <version_number>

# Then retry
./scripts/migrate.sh up
```

### "DATABASE_URL not set" Error
- Verify `.env` exists in `backend/` directory
- Ensure `DATABASE_URL` is properly formatted:
  ```
  DATABASE_URL=postgresql://user:password@localhost:5432/dbname?sslmode=disable
  ```

### Migration Failed Halfway
1. Check which version is marked dirty: `./scripts/migrate.sh version`
2. Manually fix the database state if needed
3. Force to correct version: `./scripts/migrate.sh force <version>`
4. Continue: `./scripts/migrate.sh up`

### Can't Rollback a Migration
- Verify `.down.sql` has proper reversal logic
- Check database logs for constraint violations
- May need to manually delete data before structural rollback

---

## Environment Setup

### Required .env Variables
```bash
DATABASE_URL=postgresql://user:password@host:port/database?sslmode=disable
```

### Development Database Reset (Fresh Start)
```bash
# Drop everything
./scripts/migrate.sh drop

# Reapply all migrations
./scripts/migrate.sh up
```

---

## Additional Resources

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Data Types](https://www.postgresql.org/docs/current/datatype.html)
- [Migration Best Practices](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md)

---

## Quick Reference

| Task | Command |
|------|---------|
| Create migration | `./scripts/create-migration.sh <name>` |
| Apply all migrations | `./scripts/migrate.sh up` |
| Rollback last migration | `./scripts/migrate.sh down 1` |
| Check current version | `./scripts/migrate.sh version` |
| Force to version | `./scripts/migrate.sh force <version>` |
| Drop all tables | `./scripts/migrate.sh drop` |