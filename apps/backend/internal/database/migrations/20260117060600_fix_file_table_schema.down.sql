-- Revert to original schema (if needed)
DROP TABLE IF EXISTS file;

CREATE TABLE file (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_table TEXT NOT NULL,
    parent_id UUID NOT NULL,
    role TEXT NOT NULL,
    key TEXT NOT NULL,
    file_name TEXT NOT NULL,
    file_type TEXT NOT NULL,
    file_url TEXT NOT NULL,
    content_disposition TEXT,
    polling_jwt TEXT,
    polling_url TEXT,
    custom_id TEXT,
    url TEXT NOT NULL,
    fields JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_file_parent_role ON file(parent_table, parent_id, role);