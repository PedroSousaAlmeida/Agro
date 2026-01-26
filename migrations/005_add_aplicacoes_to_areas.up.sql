ALTER TABLE areas_monitoramento
ADD COLUMN aplicacoes JSONB NOT NULL DEFAULT '[]',
ADD COLUMN updated_at TIMESTAMP;
