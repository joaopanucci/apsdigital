-- +goose Up
CREATE TABLE professions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert default professions
INSERT INTO professions (name, description) VALUES
('Enfermeiro', 'Profissional de enfermagem de nível superior'),
('Médico', 'Profissional médico'),
('Agente Comunitário de Saúde', 'Agente comunitário de saúde'),
('Odontólogo', 'Profissional de odontologia'),
('Técnico de Enfermagem', 'Técnico em enfermagem'),
('Fisioterapeuta', 'Profissional de fisioterapia'),
('Psicólogo', 'Profissional de psicologia'),
('Nutricionista', 'Profissional de nutrição'),
('Farmacêutico', 'Profissional farmacêutico'),
('Assistente Social', 'Profissional de serviço social');

-- +goose Down
DROP TABLE IF EXISTS professions;