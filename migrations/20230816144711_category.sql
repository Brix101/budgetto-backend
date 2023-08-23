-- +goose Up
-- +goose StatementBegin
CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
	name VARCHAR NOT NULL,
  note TEXT,
  user_id integer REFERENCES users (id) ON DELETE CASCADE DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE categories
-- +goose StatementEnd
