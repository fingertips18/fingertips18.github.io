CREATE TYPE skill_category AS ENUM ('frontend', 'backend', 'tools', 'others');

ALTER TABLE skill
ALTER COLUMN category TYPE skill_category
USING category::skill_category;

ALTER TABLE skill
ALTER COLUMN category SET NOT NULL;