CREATE TABLE IF NOT EXISTS users
(
    id            UUID PRIMARY KEY             DEFAULT gen_random_uuid(),
    login         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    created_at    TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS users_login ON users (login);