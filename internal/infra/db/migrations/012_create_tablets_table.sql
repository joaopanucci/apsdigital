-- +goose Up
CREATE TABLE tablets (
    id SERIAL PRIMARY KEY,
    assigned_to INTEGER REFERENCES users(id) ON DELETE SET NULL,
    municipality_id INTEGER REFERENCES municipalities(id) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'ativo',
    asset_code VARCHAR(255) UNIQUE,
    model VARCHAR(255),
    serial_number VARCHAR(255),
    assigned_at TIMESTAMP WITH TIME ZONE,
    returned_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT check_tablet_status CHECK (status IN ('ativo', 'devolvido', 'quebrado', 'furtado', 'manutencao'))
);

-- Create indexes
CREATE INDEX idx_tablets_assigned_to ON tablets(assigned_to);
CREATE INDEX idx_tablets_municipality_id ON tablets(municipality_id);
CREATE INDEX idx_tablets_status ON tablets(status);
CREATE INDEX idx_tablets_asset_code ON tablets(asset_code);

-- +goose Down
DROP TABLE IF EXISTS tablets;