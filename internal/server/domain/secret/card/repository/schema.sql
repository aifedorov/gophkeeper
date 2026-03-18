CREATE TABLE IF NOT EXISTS cards
(
    id                UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    user_id           UUID      NOT NULL REFERENCES users (id),
    name              VARCHAR(255) NOT NULL,
    encrypted_number     BYTEA     NOT NULL,
    encrypted_expired_date      BYTEA     NOT NULL,
    expired_card_holder_name BYTEA     NOT NULL,
    encrypted_cvv     BYTEA     NOT NULL,
    encrypted_notes   BYTEA,
    deleted_at        TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version           BIGINT    NOT NULL DEFAULT 1
);

CREATE UNIQUE INDEX IF NOT EXISTS cards_unique_name
    ON cards (name, user_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS name_cards ON cards (name);