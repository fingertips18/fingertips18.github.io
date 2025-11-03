-- Create the enum type
CREATE TYPE skill_category AS ENUM ('frontend', 'backend', 'tools', 'others');

-- Drop the existing default constraint
ALTER TABLE skill
ALTER COLUMN category DROP DEFAULT;

-- Convert the column type
ALTER TABLE skill
ALTER COLUMN category TYPE skill_category
USING category::skill_category;

-- NOT NULL is already set from the original column definition
-- Optionally, set a new default using an enum value (uncomment if needed):
-- ALTER TABLE skill
-- ALTER COLUMN category SET DEFAULT 'others'::skill_category;