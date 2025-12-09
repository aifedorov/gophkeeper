-- name: GetCredential :one
SELECT *
FROM credentials
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL;

-- name: CreateCredential :one
INSERT INTO credentials (user_id, name, login, password, metadata)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListCredentials :many
SELECT *
FROM credentials
WHERE user_id = $1
  AND deleted_at IS NULL;

-- name: UpdateCredential :one
UPDATE credentials
SET name       = $3,
    login      = $4,
    password   = $5,
    metadata   = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteCredential :exec
UPDATE credentials
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL;