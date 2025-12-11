-- name: GetUser :one
SELECT *
FROM users
WHERE login = $1;

-- name: CreateUser :one
INSERT INTO users (login, password_hash, salt)
VALUES ($1, $2, $3)
RETURNING *;