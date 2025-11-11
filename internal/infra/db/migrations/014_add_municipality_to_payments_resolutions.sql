-- +goose Up
-- Update payments table to include municipality_id and improve structure
ALTER TABLE payments ADD COLUMN municipality_id INTEGER REFERENCES municipalities(id);
ALTER TABLE payments ADD COLUMN competence_date DATE; -- For better date handling
ALTER TABLE payments ALTER COLUMN competence_month TYPE VARCHAR(20); -- Allow more flexible formats

-- Create index for municipality_id
CREATE INDEX idx_payments_municipality_id ON payments(municipality_id);
CREATE INDEX idx_payments_competence_date ON payments(competence_date);

-- Update resolutions table to include municipality_id and improve structure  
ALTER TABLE resolutions ADD COLUMN municipality_id INTEGER REFERENCES municipalities(id);
ALTER TABLE resolutions ADD COLUMN competence VARCHAR(20); -- Add competence field

-- Create index for municipality_id
CREATE INDEX idx_resolutions_municipality_id ON resolutions(municipality_id);
CREATE INDEX idx_resolutions_competence ON resolutions(competence);

-- +goose Down
-- Remove added columns and indexes
DROP INDEX IF EXISTS idx_payments_municipality_id;
DROP INDEX IF EXISTS idx_payments_competence_date;
DROP INDEX IF EXISTS idx_resolutions_municipality_id;
DROP INDEX IF EXISTS idx_resolutions_competence;

ALTER TABLE payments DROP COLUMN municipality_id;
ALTER TABLE payments DROP COLUMN competence_date;
ALTER TABLE resolutions DROP COLUMN municipality_id;
ALTER TABLE resolutions DROP COLUMN competence;