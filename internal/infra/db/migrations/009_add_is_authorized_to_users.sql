-- +goose Up
-- Add is_authorized field to users table
ALTER TABLE users ADD COLUMN is_authorized BOOLEAN DEFAULT FALSE;

-- Make CPF unique (it should already be, but ensuring)
-- No need to recreate the index as it already exists

-- Add index for is_authorized for faster queries
CREATE INDEX idx_users_is_authorized ON users(is_authorized);

-- +goose Down
-- Remove is_authorized field
ALTER TABLE users DROP COLUMN is_authorized;

-- Drop the index
DROP INDEX IF EXISTS idx_users_is_authorized;