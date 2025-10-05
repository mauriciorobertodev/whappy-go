CREATE TABLE IF NOT EXISTS files (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    mime TEXT NOT NULL,
    size INTEGER NOT NULL,
    sha256 TEXT NOT NULL,
    extension TEXT NOT NULL,
    path TEXT NOT NULL,
    url TEXT NOT NULL,
    width INTEGER,
    height INTEGER,
    duration INTEGER,
    pages INTEGER,
    instance_id TEXT NULL REFERENCES instances(id) ON DELETE SET NULL DEFAULT NULL,
    thumbnail_id TEXT NULL REFERENCES files(id) ON DELETE SET NULL DEFAULT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- DOWN
DROP TABLE IF EXISTS files;