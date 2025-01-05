-- +goose Up
CREATE TABLE games (
  id UUID PRIMARY KEY,
  title TEXT NOT NULL UNIQUE,
  copies SMALLINT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS games;