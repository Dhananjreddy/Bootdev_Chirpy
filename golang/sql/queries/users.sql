-- name: CreateUser :one
INSERT INTO users (id, hashed_password, created_at, updated_at, email, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    $1,
    NOW(),
    NOW(),
    $2,
    false
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

-- name: MakeChirpyRed :one
UPDATE users
SET is_chirpy_red = TRUE,  updated_at = NOW()
WHERE id = $1
RETURNING *;