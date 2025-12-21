CREATE TABLE IF NOT EXISTS files
(
    id          UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users (id),
    name            VARCHAR(255) NOT NULL,
    encrypted_path  BYTEA        NOT NULL,
    encrypted_size  BYTEA        NOT NULL,
    encrypted_notes BYTEA,
    deleted_at        TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS files_unique_name
    ON files (name, user_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_files_user_id ON files (user_id, updated_at DESC);