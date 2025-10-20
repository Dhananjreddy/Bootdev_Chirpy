-- name: CreateUser :one
INSERT INTO users (id, hashed_password, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    $1,
    NOW(),
    NOW(),
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdatePassword :one
UPDATE users
SET 
    hashed_password = $1,
    email = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING *;