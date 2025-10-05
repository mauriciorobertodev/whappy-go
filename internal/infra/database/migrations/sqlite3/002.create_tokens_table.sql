CREATE TABLE IF NOT EXISTS tokens (
    id TEXT PRIMARY KEY,
    instance_id TEXT NOT NULL REFERENCES instances(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- DOWN
DROP TABLE IF EXISTS tokens;