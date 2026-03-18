-- name: CreateCard :one
INSERT INTO cards (
                   id,
                   user_id,
                   name,
                   encrypted_number,
                   encrypted_expired_date,
                   expired_card_holder_name,
                   encrypted_cvv,
                   encrypted_notes
)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
       )
RETURNING *;

-- name: ListCards :many
SELECT *
FROM cards
WHERE user_id = $1
  AND deleted_at IS NULL;

-- name: GetCardForUpdate :one
SELECT *
FROM cards
WHERE id = $1
  AND user_id = $2
  FOR UPDATE;

-- name: UpdateCard :one
UPDATE cards
SET
    name = $4,
    encrypted_number = $5,
    encrypted_expired_date = $6,
    expired_card_holder_name = $7,
    encrypted_cvv = $8,
    encrypted_notes = $9,
    version = version + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND user_id = $2
  AND version = $3
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteCard :execrows
UPDATE cards
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL;