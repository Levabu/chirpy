-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByRefreshToken :one
SELECT * FROM users
JOIN refresh_tokens r on users.id = r.user_id
WHERE r.token = $1;

-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
gen_random_uuid(),
NOW(),
NOW(),
$1, -- email
$2  -- hashed_password
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpdateUser :one
UPDATE users
SET email = $1,
    hashed_password = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpgradeUser :one
UPDATE users
SET is_chirpy_red = TRUE,
    updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, is_chirpy_red;

-- name: DeleteUsers :exec
DELETE FROM users;