-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    cpf VARCHAR(11) NOT NULL UNIQUE,
    phone VARCHAR(20),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    profession_id UUID NOT NULL REFERENCES professions(id) ON DELETE RESTRICT,
    municipality VARCHAR(100) NOT NULL,
    unit VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'pending_authorization',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT check_status CHECK (status IN ('active', 'pending_authorization', 'blocked', 'inactive'))
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_cpf ON users(cpf);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_profession_id ON users(profession_id);
CREATE INDEX idx_users_municipality ON users(municipality);

-- +goose Down
DROP TABLE IF EXISTS users;