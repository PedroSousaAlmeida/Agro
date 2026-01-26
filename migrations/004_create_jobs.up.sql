CREATE TABLE jobs (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type            VARCHAR(50) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    payload         JSONB NOT NULL DEFAULT '{}',
    result          JSONB,
    progress        INT NOT NULL DEFAULT 0,
    total_items     INT NOT NULL DEFAULT 0,
    processed_items INT NOT NULL DEFAULT 0,
    error_count     INT NOT NULL DEFAULT 0,
    error_details   JSONB,
    started_at      TIMESTAMP,
    completed_at    TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_type ON jobs(type);
CREATE INDEX idx_jobs_created_at ON jobs(created_at DESC);
