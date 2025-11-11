-- +goose Up
CREATE TABLE tablet_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    justification TEXT,
    description TEXT,
    photos JSONB DEFAULT '[]'::jsonb,
    document_url VARCHAR(500),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP WITH TIME ZONE,
    rejected_by UUID REFERENCES users(id),
    rejected_at TIMESTAMP WITH TIME ZONE,
    rejection_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT check_type CHECK (type IN ('novo', 'devolucao', 'quebra', 'furto')),
    CONSTRAINT check_status CHECK (status IN ('pending', 'approved', 'rejected', 'completed'))
);

CREATE INDEX idx_tablet_requests_user_id ON tablet_requests(user_id);
CREATE INDEX idx_tablet_requests_type ON tablet_requests(type);
CREATE INDEX idx_tablet_requests_status ON tablet_requests(status);
CREATE INDEX idx_tablet_requests_created_at ON tablet_requests(created_at);

-- +goose Down
DROP TABLE IF EXISTS tablet_requests;