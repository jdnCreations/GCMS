-- +goose Up
CREATE TABLE genres (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE game_genres (
  game_id UUID REFERENCES games(id) ON DELETE CASCADE,
  genre_id UUID REFERENCES genres(id) ON DELETE CASCADE,
  PRIMARY KEY (game_id, genre_id)
);

-- +goose Down
DROP TABLE IF EXISTS game_genres;
DROP TABLE IF EXISTS genres;