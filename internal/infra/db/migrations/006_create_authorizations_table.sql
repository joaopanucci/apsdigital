-- +goose Up
CREATE TABLE authorizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    authorized_by UUID REFERENCES users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    comments TEXT,
    authorized_at TIMESTAMP WITH TIME ZONE,
    rejected_at TIMESTAMP WITH TIME ZONE,
    rejection_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT check_status CHECK (status IN ('pending', 'approved', 'rejected'))
);

CREATE INDEX idx_authorizations_user_id ON authorizations(user_id);
CREATE INDEX idx_authorizations_status ON authorizations(status);
CREATE INDEX idx_authorizations_created_at ON authorizations(created_at);

-- +goose Down
DROP TABLE IF EXISTS authorizations;