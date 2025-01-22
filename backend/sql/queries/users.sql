-- name: GetUserById :one
SELECT id, first_name, last_name, email, is_admin FROM users where id = $1;

-- name: GetUserByEmail :one
SELECT * from users where email = $1;

-- name: CreateUser :one
INSERT INTO users (id, first_name, last_name, email, hashed_password)
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3,
  $4
)
RETURNING id, first_name, last_name, email, is_admin;

-- name: DeleteUserById :exec
DELETE FROM users where id = $1;

-- name: GetAllUsers :many
SELECT id, first_name, last_name, email, is_admin FROM users ORDER BY first_name;

-- name: UpdateUser :one
UPDATE users 
SET first_name = COALESCE(NULLIF($1, ''), first_name),
    last_name = COALESCE(NULLIF($2, ''), last_name),
    email = COALESCE(NULLIF($3, ''), email)
WHERE id = $4
RETURNING id, first_name, last_name, email;

-- name: SetAdmin :exec
UPDATE users
SET is_admin = $1
WHERE id = $2;

-- name: GetUserFromRefreshToken :one
SELECT *
FROM users 
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
where refresh_tokens.token = $1;

-- name: GetUserStats :one
SELECT 
  (SELECT COUNT(*) FROM users) as total_users,
  (SELECT COUNT(DISTINCT users.id) FROM users
   JOIN reservations on users.id = reservations.user_id) as users_with_reservations;
