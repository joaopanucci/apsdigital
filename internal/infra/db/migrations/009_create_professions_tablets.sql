-- Criar tabela de profissões da área da saúde
-- Migration: 009_create_professions_tablets.sql

-- Criar tabela de profissões
CREATE TABLE IF NOT EXISTS professions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Inserir profissões oficiais da área da saúde
INSERT INTO professions (name) VALUES
    ('Enfermeiro'),
    ('Médico'),
    ('Técnico de Enfermagem'),
    ('Odontólogo'),
    ('Auxiliar Odontológico'),
    ('ACS - Agente Comunitário de Saúde'),
    ('Psicólogo'),
    ('Assistente Social'),
    ('Nutricionista'),
    ('Fisioterapeuta'),
    ('Farmacêutico'),
    ('Terapeuta Ocupacional'),
    ('Fonoaudiólogo'),
    ('Sanitarista'),
    ('Biomédico'),
    ('Educador Físico')
ON CONFLICT (name) DO NOTHING;

-- Criar tabela de tablets
CREATE TABLE IF NOT EXISTS tablets (
    id SERIAL PRIMARY KEY,
    assigned_to UUID REFERENCES users(id),
    status VARCHAR(50) NOT NULL DEFAULT 'ativo', -- ativo, devolvido, quebrado, furtado
    asset_code VARCHAR(255) UNIQUE,
    municipality_id INTEGER REFERENCES municipalities(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar índices para performance
CREATE INDEX IF NOT EXISTS idx_tablets_assigned_to ON tablets(assigned_to);
CREATE INDEX IF NOT EXISTS idx_tablets_municipality_id ON tablets(municipality_id);
CREATE INDEX IF NOT EXISTS idx_tablets_status ON tablets(status);

-- Comentários
COMMENT ON TABLE professions IS 'Profissões da área da saúde';
COMMENT ON TABLE tablets IS 'Tablets distribuídos aos agentes';
COMMENT ON COLUMN tablets.status IS 'Status do tablet: ativo, devolvido, quebrado, furtado';
COMMENT ON COLUMN tablets.asset_code IS 'Código patrimonial do tablet';