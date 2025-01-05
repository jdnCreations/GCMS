-- name: GetGenreById :one
SELECT * FROM genres where id = $1;

-- name: CreateGenre :one
INSERT INTO genres (id, name)
VALUES (
  gen_random_uuid(),
  $1
)
RETURNING *;

-- name: DeleteGenreById :exec
DELETE FROM genres WHERE id = $1;

-- name: GetAllGenres :many
SELECT * from genres;

-- name: UpdateGenrename :one
UPDATE genres 
SET name = $1
WHERE id = $2
RETURNING id, name;


-- name: UpdateGenreCopies :one
UPDATE genres 
SET name = $1
WHERE id = $2
RETURNING id, name;


-- name: UpdateGenre :one
UPDATE genres 
SET name = COALESCE(NULLIF($1, ''), name)
WHERE id = $2
RETURNING id, name; 