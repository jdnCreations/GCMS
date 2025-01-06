-- name: AddGenreToGame :one
INSERT INTO game_genres (game_id, genre_id)
VALUES (
  $1,
  $2
)
RETURNING *;