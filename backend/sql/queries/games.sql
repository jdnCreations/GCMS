-- name: GetGameById :one
SELECT * FROM games where id = $1;

-- name: CreateGame :one
INSERT INTO games (id, title, copies)
VALUES (
  gen_random_uuid(),
  $1,
  $2
)
RETURNING *;

-- name: DeleteGameById :exec
DELETE FROM games WHERE id = $1;

/* maybe add functionality to pass in what to sort by? */
-- name: GetAllGames :many
SELECT * from games;

-- name: UpdateGameTitle :one
UPDATE games 
SET title = $1
WHERE id = $2
RETURNING id, title, copies;


-- name: UpdateGameCopies :one
UPDATE games 
SET copies = $1
WHERE id = $2
RETURNING id, title, copies;


-- name: UpdateGame :one
UPDATE games 
SET title = COALESCE(NULLIF($1, ''), title),
    copies = CASE WHEN $2::SMALLINT IS NOT NULL THEN $2 ELSE copies END
WHERE id = $3
RETURNING id, title, copies;