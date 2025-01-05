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


-- name: UpdateGame :one
UPDATE games
SET
  title = CASE
    WHEN $1::TEXT IS NOT NULL THEN $1
    ELSE title
  END,
  copies = CASE
    WHEN $2::SMALLINT IS NOT NULL then $2
    ELSE copies
  END
WHERE id = $3 RETURNING id, title, copies;