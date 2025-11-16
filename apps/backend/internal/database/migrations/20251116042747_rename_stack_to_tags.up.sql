-- Rename stack column to tags
ALTER TABLE project
RENAME COLUMN stack TO tags;
