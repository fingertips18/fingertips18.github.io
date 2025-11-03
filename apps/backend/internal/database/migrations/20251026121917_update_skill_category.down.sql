ALTER TABLE skill
ALTER COLUMN category DROP NOT NULL;

ALTER TABLE skill
ALTER COLUMN category TYPE TEXT
USING category::text;

DROP TYPE IF EXISTS skill_category;
