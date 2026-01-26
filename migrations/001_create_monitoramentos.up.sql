CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE monitoramentos (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    data_upload   TIMESTAMP NOT NULL DEFAULT NOW(),
    nome_arquivo  VARCHAR(255) NOT NULL,
    status        VARCHAR(20) NOT NULL CHECK (status IN ('processando', 'concluido', 'erro')),
    total_linhas  INT NOT NULL DEFAULT 0,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_monitoramentos_status ON monitoramentos(status);
CREATE INDEX idx_monitoramentos_data_upload ON monitoramentos(data_upload DESC);
