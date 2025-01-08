-- name: CreateReservation :one
INSERT INTO reservations (
  id, 
  start_time, 
  end_time, 
  user_id, 
  game_id
  )
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3,
  $4
)
RETURNING *;

-- name: GetAllReservations :many
SELECT * from reservations;

-- name: GetReservationsForUser :many
SELECT * from reservations where user_id = $1;

-- name: GetAllActiveReservations :many
SELECT * from reservations where end_time > NOW();