-- +goose Up
CREATE TABLE resolutions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    year INTEGER NOT NULL,
    type VARCHAR(10) NOT NULL,
    description TEXT,
    original_file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    uploaded_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT check_type CHECK (type IN ('MS', 'SES'))
);

CREATE INDEX idx_resolutions_year ON resolutions(year);
CREATE INDEX idx_resolutions_type ON resolutions(type);
CREATE INDEX idx_resolutions_uploaded_by ON resolutions(uploaded_by);
CREATE INDEX idx_resolutions_created_at ON resolutions(created_at);

-- +goose Down
DROP TABLE IF EXISTS resolutions;