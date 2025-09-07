CREATE TABLE project (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    preview TEXT NOT NULL,
    blur_hash TEXT,
    title TEXT NOT NULL,
    sub_title TEXT,
    description TEXT,
    stack JSONB,
    type TEXT,
    link TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);