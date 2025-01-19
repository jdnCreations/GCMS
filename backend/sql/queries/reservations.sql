-- name: CreateReservation :one
INSERT INTO reservations (
  id, 
  res_date,
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
  $4,
  $5
)
RETURNING *;

-- name: GetAllReservations :many
SELECT * from reservations;

-- name: GetReservationsForUser :many
SELECT reservations.*,
       games.title as game_name
FROM reservations
JOIN
  games
ON
  reservations.game_id = games.id
where user_id = $1;

-- name: GetAllActiveReservations :many
SELECT * from reservations where end_time > NOW();

-- name: GetExpiredReservations :many
SELECT * from reservations where end_time < NOW()::TIME AND active = true;

-- name: SetReservationInactive :one
UPDATE reservations
SET active = false
WHERE id = $1 AND
end_time < NOW()::TIME
RETURNING *;

-- name: CheckGameReservation :one
SELECT
  COALESCE((SELECT copies FROM games WHERE game_id = $1), 0) - COALESCE(COUNT(*), 0) AS available_copies
FROM reservations
WHERE reservations.game_id = $1
AND reservations.start_time < $2
AND reservations.end_time > $3
GROUP BY reservations.game_id;

-- name: DeleteReservation :exec
DELETE FROM reservations where id = $1;

-- name: GetReservationById :one
SELECT * from reservations where id = $1;