-- Alterações para implementar login via CPF e autorização de usuários
-- Migration: 008_update_users_cpf_auth.sql

-- Adicionar índice único no CPF
CREATE UNIQUE INDEX IF NOT EXISTS users_cpf_unique ON users(cpf) WHERE cpf IS NOT NULL AND cpf != '';

-- Adicionar campo is_authorized se não existir
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'is_authorized') THEN
        ALTER TABLE users ADD COLUMN is_authorized BOOLEAN DEFAULT false;
    END IF;
END $$;

-- Atualizar usuários existentes para serem autorizados (para não quebrar o sistema atual)
UPDATE users SET is_authorized = true WHERE is_authorized IS NULL;

-- Comentários
COMMENT ON COLUMN users.cpf IS 'CPF do usuário - usado como chave principal de login';
COMMENT ON COLUMN users.is_authorized IS 'Define se o usuário foi autorizado a usar o sistema';