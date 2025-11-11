-- +goose Up
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file_url VARCHAR(500) NOT NULL,
    competence_month VARCHAR(7) NOT NULL, -- YYYY-MM format
    competence_year INTEGER NOT NULL,
    municipality VARCHAR(100) NOT NULL,
    uploaded_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    original_file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_payments_competence_month ON payments(competence_month);
CREATE INDEX idx_payments_competence_year ON payments(competence_year);
CREATE INDEX idx_payments_municipality ON payments(municipality);
CREATE INDEX idx_payments_uploaded_by ON payments(uploaded_by);
CREATE INDEX idx_payments_created_at ON payments(created_at);

-- +goose Down
DROP TABLE IF EXISTS payments;