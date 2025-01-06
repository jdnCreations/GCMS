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



/* -- +goose Up
CREATE TABLE reservations (
  id UUID PRIMARY KEY,
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP NOT NULL,
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  game_id UUID REFERENCES games(id) ON DELETE CASCADE,
  CHECK (start_time < end_time)
);
 */
