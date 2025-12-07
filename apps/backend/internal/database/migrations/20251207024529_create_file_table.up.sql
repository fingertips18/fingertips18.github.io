CREATE TABLE file {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_table TEXT NOT NULL,    -- e.g., 'project', 'user', etc.
    parent_id UUID NOT NULL,       -- the id of the parent row
    role TEXT NOT NULL,            -- e.g., 'preview', 'attachment'
    key TEXT NOT NULL,
    file_name TEXT NOT NULL,
    file_url TEXT NOT NULL,
    content_disposition TEXT NOT NULL,
    polling_jwt TEXT NOT NULL,
    polling_url TEXT NOT NULL,
    custom_id TEXT NOT NULL,
    url TEXT NOT NULL,
    fields JSONB DEFAULT {},
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
}

-- Index for fast lookup
CREATE INDEX idx_file_parent_role ON file(parent_table, parent_id, role);