CREATE TABLE IF NOT EXISTS files (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mime VARCHAR(255) NOT NULL,
    size INTEGER NOT NULL,
    sha256 VARCHAR(64) NOT NULL,
    extension VARCHAR(10) NOT NULL,
    path TEXT NOT NULL,
    url TEXT NOT NULL,
    width INTEGER,
    height INTEGER,
    duration INTEGER,
    pages INTEGER,
    instance_id VARCHAR(36) NULL REFERENCES instances(id) ON DELETE SET NULL,
    thumbnail_id VARCHAR(36) NULL REFERENCES files(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

-- DOWN
DROP TABLE IF EXISTS files;