-- Down migration
ALTER TABLE project
RENAME COLUMN tags TO stack;