CREATE TABLE IF NOT EXISTS files
(
    id          UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users (id),
    filename    VARCHAR(255) NOT NULL,
    file_path   TEXT         NOT NULL,
    file_size   BIGINT       NOT NULL,
    mime_type   VARCHAR(255) NOT NULL,
    uploaded_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS files_unique_name
    ON files (filename, user_id);

CREATE INDEX IF NOT EXISTS idx_files_user_id ON files (user_id, uploaded_at DESC);