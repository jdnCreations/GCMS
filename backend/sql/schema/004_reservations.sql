-- +goose Up
CREATE TABLE reservations (
  id UUID PRIMARY KEY,
  res_date DATE NOT NULL,
  start_time TIME NOT NULL,
  end_time TIME NOT NULL,
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  game_id UUID REFERENCES games(id) ON DELETE CASCADE,
  CHECK (start_time < end_time)
);

-- +goose Down
DROP TABLE IF EXISTS reservations;