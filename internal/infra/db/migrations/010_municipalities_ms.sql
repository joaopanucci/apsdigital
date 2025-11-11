-- Adicionar municípios de Mato Grosso do Sul com códigos IBGE
-- Migration: 010_municipalities_ms.sql

-- Atualizar estrutura da tabela municipalities para incluir código IBGE
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'municipalities' AND column_name = 'ibge_code') THEN
        ALTER TABLE municipalities ADD COLUMN ibge_code VARCHAR(7);
    END IF;
END $$;

-- Limpar dados existentes para inserir os municípios corretos do MS
DELETE FROM municipalities;

-- Inserir os 79 municípios de Mato Grosso do Sul
INSERT INTO municipalities (name, ibge_code) VALUES
    ('Água Clara', '5000203'),
    ('Alcinópolis', '5000252'),
    ('Amambai', '5000609'),
    ('Anastácio', '5000708'),
    ('Anaurilândia', '5000807'),
    ('Angélica', '5000856'),
    ('Antônio João', '5000906'),
    ('Aparecida do Taboado', '5001003'),
    ('Aquidauana', '5001102'),
    ('Aral Moreira', '5001243'),
    ('Bandeirantes', '5001508'),
    ('Bataguassu', '5001904'),
    ('Batayporã', '5001953'),
    ('Bela Vista', '5002001'),
    ('Bodoquena', '5002100'),
    ('Bonito', '5002159'),
    ('Brasilândia', '5002209'),
    ('Caarapó', '5002308'),
    ('Camapuã', '5002407'),
    ('Campo Grande', '5002704'),
    ('Caracol', '5002803'),
    ('Cassilândia', '5002902'),
    ('Chapadão do Sul', '5003108'),
    ('Corguinho', '5003157'),
    ('Coronel Sapucaia', '5003207'),
    ('Corumbá', '5003256'),
    ('Costa Rica', '5003306'),
    ('Coxim', '5003454'),
    ('Deodápolis', '5003488'),
    ('Dois Irmãos do Buriti', '5003504'),
    ('Douradina', '5003702'),
    ('Dourados', '5003751'),
    ('Eldorado', '5003801'),
    ('Fátima do Sul', '5003900'),
    ('Figueirão', '5004007'),
    ('Glória de Dourados', '5004106'),
    ('Guia Lopes da Laguna', '5004304'),
    ('Iguatemi', '5004403'),
    ('Inocência', '5004502'),
    ('Itaporã', '5004601'),
    ('Itaquiraí', '5004700'),
    ('Ivinhema', '5004809'),
    ('Japorã', '5004908'),
    ('Jaraguari', '5005004'),
    ('Jardim', '5005103'),
    ('Jateí', '5005152'),
    ('Juti', '5005202'),
    ('Ladário', '5005251'),
    ('Laguna Carapã', '5005400'),
    ('Maracaju', '5005608'),
    ('Miranda', '5005681'),
    ('Mundo Novo', '5005707'),
    ('Naviraí', '5005806'),
    ('Nioaque', '5005905'),
    ('Nova Alvorada do Sul', '5006002'),
    ('Nova Andradina', '5006101'),
    ('Novo Horizonte do Sul', '5006200'),
    ('Paranaíba', '5006259'),
    ('Paranhos', '5006309'),
    ('Pedro Gomes', '5006358'),
    ('Ponta Porã', '5006606'),
    ('Porto Murtinho', '5006903'),
    ('Ribas do Rio Pardo', '5007109'),
    ('Rio Brilhante', '5007208'),
    ('Rio Negro', '5007307'),
    ('Rio Verde de Mato Grosso', '5007406'),
    ('Rochedo', '5007505'),
    ('Santa Rita do Pardo', '5007554'),
    ('São Gabriel do Oeste', '5007695'),
    ('Selvíria', '5007802'),
    ('Sete Quedas', '5007901'),
    ('Sidrolândia', '5008008'),
    ('Sonora', '5008305'),
    ('Tacuru', '5008404'),
    ('Taquarussu', '5008503'),
    ('Terenos', '5008602'),
    ('Três Lagoas', '5008701'),
    ('Vicentina', '5008800')
ON CONFLICT (name) DO NOTHING;

-- Criar índices
CREATE INDEX IF NOT EXISTS idx_municipalities_ibge_code ON municipalities(ibge_code);

-- Comentários
COMMENT ON COLUMN municipalities.ibge_code IS 'Código IBGE do município';