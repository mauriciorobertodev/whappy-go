CREATE TABLE IF NOT EXISTS webhooks (
    id TEXT PRIMARY KEY,
    url TEXT NOT NULL,
    events TEXT NOT NULL,
    secret TEXT,
    active BOOLEAN NOT NULL DEFAULT FALSE,
    
    updated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,

    instance_id TEXT NOT NULL REFERENCES instances(id) ON DELETE CASCADE
);

-- DOWN
DROP TABLE IF EXISTS webhooks;