-- +goose Up
-- Add municipality_id to users table
ALTER TABLE users ADD COLUMN municipality_id INTEGER REFERENCES municipalities(id);

-- Create index for faster queries
CREATE INDEX idx_users_municipality_id ON users(municipality_id);

-- Remove the old municipality VARCHAR field (will be replaced by municipality_id)
-- We'll keep it for now to avoid breaking changes, but it should be phased out

-- +goose Down
-- Remove municipality_id field and index
DROP INDEX IF EXISTS idx_users_municipality_id;
ALTER TABLE users DROP COLUMN municipality_id;