CREATE TABLE IF NOT EXISTS skill (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    icon TEXT NOT NULL,
    hex_color TEXT NOT NULL CHECK (hex_color ~ '^#[0-9A-Fa-f]{6}$'),
    label TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create function to automatically update updated_at
CREATE OR REPLACE FUNCTION update_skill_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
CREATE TRIGGER trigger_update_skill_updated_at
    BEFORE UPDATE ON skill
    FOR EACH ROW
    EXECUTE FUNCTION update_skill_updated_at();

-- Add indexes for query performance
CREATE INDEX idx_skill_label ON skill(label);