ALTER TABLE skill ADD COLUMN category TEXT NOT NULL DEFAULT '';

-- Add index for category queries
CREATE INDEX idx_skill_category ON skill(category);