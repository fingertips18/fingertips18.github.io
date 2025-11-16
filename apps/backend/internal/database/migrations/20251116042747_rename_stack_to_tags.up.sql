-- Rename stack column to tags
ALTER TABLE project
RENAME COLUMN stack TO tags;

-- Ensure NOT NULL and default are preserved
ALTER TABLE project
ALTER COLUMN tags SET NOT NULL,
ALTER COLUMN tags SET DEFAULT '{}';
