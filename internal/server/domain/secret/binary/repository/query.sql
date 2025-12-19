-- name: CreateFile :exec
INSERT INTO files (id, user_id, filename, file_path, file_size, mime_type, uploaded_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: ListFiles :many
SELECT *
FROM files
WHERE user_id = $1
ORDER BY uploaded_at DESC;

-- name: DeleteFile :execrows
DELETE
FROM files
WHERE id = $1
  AND user_id = $2;