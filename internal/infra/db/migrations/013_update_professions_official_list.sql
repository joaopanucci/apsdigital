-- +goose Up
-- Update professions table with official health professions
DELETE FROM professions; -- Clear existing data

INSERT INTO professions (name, description, is_active) VALUES
('Enfermeiro', 'Profissional de enfermagem de nível superior', true),
('Médico', 'Profissional médico', true),
('Técnico de Enfermagem', 'Profissional técnico de enfermagem', true),
('Odontólogo', 'Profissional de odontologia (cirurgião-dentista)', true),
('Auxiliar Odontológico', 'Auxiliar em saúde bucal', true),
('Agente Comunitário de Saúde', 'Agente comunitário de saúde (ACS)', true),
('Psicólogo', 'Profissional de psicologia', true),
('Assistente Social', 'Profissional de serviço social', true),
('Nutricionista', 'Profissional de nutrição', true),
('Fisioterapeuta', 'Profissional de fisioterapia', true),
('Farmacêutico', 'Profissional farmacêutico', true),
('Terapeuta Ocupacional', 'Profissional de terapia ocupacional', true),
('Fonoaudiólogo', 'Profissional de fonoaudiologia', true),
('Sanitarista', 'Profissional sanitarista', true),
('Biomédico', 'Profissional biomédico', true),
('Educador Físico', 'Profissional de educação física', true);

-- +goose Down
-- Restore previous data (this is just an example - in real scenario you'd backup first)
DELETE FROM professions;
INSERT INTO professions (name, description, is_active) VALUES
('Enfermeiro', 'Profissional de enfermagem de nível superior', true),
('Médico', 'Profissional médico', true),
('Agente Comunitário de Saúde', 'Agente comunitário de saúde', true),
('Odontólogo', 'Profissional de odontologia', true),
('Técnico de Enfermagem', 'Técnico em enfermagem', true);