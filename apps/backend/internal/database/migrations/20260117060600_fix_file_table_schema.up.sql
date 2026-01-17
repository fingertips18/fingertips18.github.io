-- Rename columns to match new schema
ALTER TABLE file RENAME COLUMN file_name TO name;
ALTER TABLE file RENAME COLUMN file_type TO type;

-- Preserve URL data (choose canonical source)
UPDATE file SET url = COALESCE(url, file_url);

-- Add size column and remove obsolete ones
ALTER TABLE file
  ADD COLUMN IF NOT EXISTS size BIGINT NOT NULL DEFAULT 0,
  DROP COLUMN IF EXISTS key,
  DROP COLUMN IF EXISTS file_url,
  DROP COLUMN IF EXISTS content_disposition,
  DROP COLUMN IF EXISTS polling_jwt,
  DROP COLUMN IF EXISTS polling_url,
  DROP COLUMN IF EXISTS custom_id,
  DROP COLUMN IF EXISTS fields;

-- Ensure index exists (idempotent)
CREATE INDEX IF NOT EXISTS idx_file_parent_role ON file(parent_table, parent_id, role);