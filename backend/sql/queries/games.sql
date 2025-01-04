-- name: GetGameById :one
SELECT * FROM games where id = $1;

-- name: CreateGame :one
INSERT INTO games (id, title, genre, copies)
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3
)
RETURNING *;

-- name: DeleteGameById :exec
DELETE FROM games WHERE id = $1;

/* maybe add functionality to pass in what to sort by? */
-- name: GetAllGames :many
SELECT * from games;

-- name: UpdateGame :one
UPDATE games 
SET title = COALESCE(NULLIF($1, ''), title),
    genre = COALESCE(NULLIF($2, ''), genre),
    copies = COALESCE(NULLIF($3, ''), copies)
WHERE id = $4
RETURNING id, title, genre, copies;

