-- name: GetUser :one
SELECT *
FROM users
WHERE login = $1;

-- name: CreateUser :one
INSERT INTO users (id, login, password_hash, salt)
VALUES ($1, $2, $3, $4)
RETURNING *;