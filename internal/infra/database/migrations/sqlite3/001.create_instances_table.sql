CREATE TABLE IF NOT EXISTS instances (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT,
    jid TEXT,
    lid TEXT,
    device TEXT,
    status TEXT NOT NULL,

    last_login_at TIMESTAMP,
    last_connected_at TIMESTAMP,
    banned_at TIMESTAMP,
    ban_expires_at TIMESTAMP,

    updated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- DOWN
DROP TABLE IF EXISTS instances;