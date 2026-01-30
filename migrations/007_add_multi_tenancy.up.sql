ALTER TABLE monitoramentos
    ADD COLUMN client_id UUID REFERENCES clients(id) ON DELETE CASCADE,
    ADD COLUMN user_id VARCHAR(100);

CREATE INDEX idx_monitoramentos_client_id ON monitoramentos(client_id);
CREATE INDEX idx_monitoramentos_user_id ON monitoramentos(user_id);

ALTER TABLE areas_monitoramento
    ADD COLUMN client_id UUID REFERENCES clients(id) ON DELETE CASCADE,
    ADD COLUMN user_id VARCHAR(100);

CREATE INDEX idx_areas_monitoramento_client_id ON areas_monitoramento(client_id);
CREATE INDEX idx_areas_monitoramento_user_id ON areas_monitoramento(user_id);

ALTER TABLE jobs
    ADD COLUMN client_id UUID REFERENCES clients(id) ON DELETE CASCADE,
    ADD COLUMN user_id VARCHAR(100);

CREATE INDEX idx_jobs_client_id ON jobs(client_id);
CREATE INDEX idx_jobs_user_id ON jobs(user_id);
