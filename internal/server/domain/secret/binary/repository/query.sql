-- name: CreateFile :one
INSERT INTO files (id, user_id, name, encrypted_path, encrypted_size, encrypted_notes, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFile :one
SELECT *
FROM files
WHERE id = $1 AND
    user_id = $2 AND
    deleted_at IS NULL;

-- name: ListFiles :many
SELECT *
FROM files
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY updated_at DESC;

-- name: GetFileForUpdate :one
SELECT *
FROM files
WHERE id = $1
  AND user_id = $2
FOR UPDATE;

-- name: UpdateFile :one
UPDATE files
SET
    name = $3,
    encrypted_path = $4,
    encrypted_size = $5,
    encrypted_notes = $6,
    version = version + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND
    user_id = $2 AND
    deleted_at IS NULL
RETURNING *;

-- name: DeleteFile :execrows
DELETE
FROM files
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL;