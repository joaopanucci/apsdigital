-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    level INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert default roles
INSERT INTO roles (name, description, level) VALUES
('ADM', 'Administrador - Controle total do sistema', 1),
('Coordenador', 'Coordenador - Pode aprovar ações de gerente e ACS', 2),
('Gerente', 'Gerente - Aprova/atua sobre ACS', 3),
('ACS', 'Agente Comunitário de Saúde - Envia solicitações', 4);

-- +goose Down
DROP TABLE IF EXISTS roles;
DROP EXTENSION IF EXISTS "uuid-ossp";