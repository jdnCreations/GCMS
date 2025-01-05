-- name: GetUserById :one
SELECT * FROM users where id = $1;

-- name: CreateUser :one
INSERT INTO users (id, first_name, last_name, email)
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3
)
RETURNING *;

-- name: DeleteUserById :exec
DELETE FROM users where id = $1;

-- name: GetAllUsers :many
SELECT * FROM users ORDER BY first_name;

-- name: UpdateUser :one
UPDATE users 
SET first_name = COALESCE(NULLIF($1, ''), first_name),
    last_name = COALESCE(NULLIF($2, ''), last_name),
    email = COALESCE(NULLIF($3, ''), email)
WHERE id = $4
RETURNING id, first_name, last_name, email;

-- name: SetAdmin :one
UPDATE users
SET is_admin = $1
WHERE id = $2
RETURNING *;