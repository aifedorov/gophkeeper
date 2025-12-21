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

-- name: UpdateCard :one
UPDATE cards
SET
    name = $3,
    encrypted_number = $4,
    encrypted_expired_date = $5,
    expired_card_holder_name = $6,
    encrypted_cvv = $7,
    encrypted_notes = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteCard :execrows
UPDATE cards
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL;