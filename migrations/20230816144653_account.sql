-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  balance DOUBLE PRECISION DEFAULT 0,
  note TEXT,
  user_id integer REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accounts
-- +goose StatementEnd
