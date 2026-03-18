CREATE TABLE IF NOT EXISTS credentials
(
    id                UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    user_id           UUID      NOT NULL REFERENCES users (id),
    name           VARCHAR(255) NOT NULL,
    encryptedLogin    BYTEA     NOT NULL,
    encryptedPassword BYTEA     NOT NULL,
    encryptedNotes BYTEA,
    deleted_at     TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version    BIGINT    NOT NULL DEFAULT 1
);

CREATE UNIQUE INDEX IF NOT EXISTS credentials_unique_name
    ON credentials (name, user_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS name_credentials ON credentials (name);