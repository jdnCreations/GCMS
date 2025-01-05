/* -- name: GetGenreById :one
SELECT * FROM genres where id = $1;

-- name: CreateGenre :one
INSERT INTO genres (id, title)
VALUES (
  gen_random_uuid(),
  $1,
)
RETURNING *;

-- name: DeleteGenreById :exec
DELETE FROM genres WHERE id = $1;

/* maybe add functionality to pass in what to sort by? */
-- name: GetAllGenres :many
SELECT * from genres;

-- name: UpdateGenreTitle :one
UPDATE genres 
SET title = $1
WHERE id = $2
RETURNING id, title, copies;


-- name: UpdateGenreCopies :one
UPDATE genres 
SET copies = $1
WHERE id = $2
RETURNING id, title, copies;


-- name: UpdateGenre :one
UPDATE genres 
SET title = COALESCE(NULLIF($1, ''), title),
    copies = CASE WHEN $2::SMALLINT IS NOT NULL THEN $2 ELSE copies END
WHERE id = $3
RETURNING id, title, copies; */