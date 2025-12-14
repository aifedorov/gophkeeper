-- name: CreateCredential :one
INSERT INTO credentials (user_id, name, encryptedLogin, encryptedPassword, encryptedNotes)
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
    encryptedLogin    = $4,
    encryptedPassword = $5,
    encryptedNotes    = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteCredential :execrows
UPDATE credentials
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL;