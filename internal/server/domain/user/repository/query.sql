-- name: GetUser :one
SELECT *
FROM users
WHERE login = $1
  AND password_hash = $2;

-- name: CreateUser :one
INSERT INTO users (login, password_hash)
VALUES ($1, $2)
RETURNING *;