CREATE TABLE IF NOT EXISTS instances (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    jid VARCHAR(50),
    lid VARCHAR(50),
    device VARCHAR(50),
    status VARCHAR(50) NOT NULL,
    last_login_at TIMESTAMPTZ,
    last_connected_at TIMESTAMPTZ,
    banned_at TIMESTAMPTZ,
    ban_expires_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

-- DOWN
DROP TABLE IF EXISTS instances;