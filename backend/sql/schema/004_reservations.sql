-- +goose Up
CREATE TABLE reservations (
  id UUID PRIMARY KEY,
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP NOT NULL,
  customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
  game_id UUID REFERENCES games(id) ON DELETE CASCADE,
  CHECK (start_time < end_time)
);

-- +goose Down
DROP TABLE IF EXISTS reservations;