CREATE TABLE IF NOT EXISTS files
(
    id          UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users (id),
    name            VARCHAR(255) NOT NULL,
    encrypted_path  BYTEA        NOT NULL,
    encrypted_size  BYTEA        NOT NULL,
    encrypted_notes BYTEA,
    uploaded_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_files_user_id ON files (user_id, uploaded_at DESC);