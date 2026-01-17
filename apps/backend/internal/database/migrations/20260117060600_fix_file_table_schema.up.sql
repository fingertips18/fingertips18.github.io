-- Drop the old table
DROP TABLE IF EXISTS file;

-- Recreate with correct schema
CREATE TABLE file (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_table TEXT NOT NULL,
    parent_id UUID NOT NULL,
    role TEXT NOT NULL,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    type TEXT NOT NULL,
    size BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Index for fast lookup
CREATE INDEX idx_file_parent_role ON file(parent_table, parent_id, role);