CREATE TABLE IF NOT EXISTS migration_history (
    id            SERIAL PRIMARY KEY,
    version       INT NOT NULL,
    name          VARCHAR(255) NOT NULL,
    direction     VARCHAR(10) NOT NULL CHECK (direction IN ('up', 'down')),
    executed_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    execution_ms  INT
);

CREATE INDEX idx_migration_history_version ON migration_history(version);

-- Insere as migrations já executadas
INSERT INTO migration_history (version, name, direction, executed_at) VALUES
(1, 'create_monitoramentos', 'up', NOW()),
(2, 'create_areas_monitoramento', 'up', NOW());
