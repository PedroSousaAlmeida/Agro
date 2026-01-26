CREATE TABLE areas_monitoramento (
    id                 UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    monitoramento_id   UUID NOT NULL REFERENCES monitoramentos(id) ON DELETE CASCADE,
    setor              VARCHAR(100),
    setor2             VARCHAR(100),
    cod_fazenda        VARCHAR(50),
    desc_fazenda       VARCHAR(255),
    quadra             VARCHAR(50),
    corte              INT,
    area_total         DECIMAL(10,2),
    desc_textura_solo  VARCHAR(100),
    corte_atual        INT,
    reforma            VARCHAR(50),
    mes_colheita       VARCHAR(20),
    restricao          VARCHAR(255),
    pragas_data        JSONB NOT NULL DEFAULT '{"pragas": {}}',
    created_at         TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_areas_monitoramento_id ON areas_monitoramento(monitoramento_id);
CREATE INDEX idx_areas_cod_fazenda ON areas_monitoramento(cod_fazenda);
CREATE INDEX idx_areas_pragas_data ON areas_monitoramento USING GIN (pragas_data);
