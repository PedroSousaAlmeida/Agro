CREATE TABLE clients (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name                VARCHAR(255) NOT NULL,
    slug                VARCHAR(100) NOT NULL UNIQUE,
    max_users           INT NOT NULL DEFAULT 10,
    active              BOOLEAN NOT NULL DEFAULT true,
    metadata            JSONB DEFAULT '{}',
    keycloak_group_id   VARCHAR(100),
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_clients_slug ON clients(slug);
CREATE INDEX idx_clients_active ON clients(active);

CREATE TABLE client_users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id       UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    user_id         VARCHAR(100) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    role            VARCHAR(50) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    active          BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(client_id, user_id)
);

CREATE INDEX idx_client_users_client_id ON client_users(client_id);
CREATE INDEX idx_client_users_user_id ON client_users(user_id);
CREATE UNIQUE INDEX idx_client_users_email_client ON client_users(client_id, email);
