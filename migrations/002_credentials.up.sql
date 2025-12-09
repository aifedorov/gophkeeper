CREATE TABLE IF NOT EXISTS credentials
(
    id         UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL REFERENCES users (id),
    name       VARCHAR(255) NOT NULL,
    login      VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    metadata   TEXT         NOT NULL,
    deleted_at TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS credentials_unique_name
    ON credentials (name, user_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_credentials_user_id ON credentials (user_id);