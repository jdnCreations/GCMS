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

-- name: CheckGameReservation :one
SELECT
  COUNT(*) as reserved_count
FROM reservations
WHERE reservations.game_id = $1
AND reservations.start_time < $2
AND reservations.end_time > $3
HAVING
  COUNT(*) < (SELECT copies from games WHERE id = $1);

-- name: DeleteReservation :exec
DELETE FROM reservations where id = $1;