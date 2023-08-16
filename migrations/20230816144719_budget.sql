-- +goose Up
-- +goose StatementBegin
-- Add migration script here
CREATE TABLE budgets (
    id SERIAL PRIMARY KEY,
    amount DOUBLE PRECISION DEFAULT 0,
    category_id integer UNIQUE REFERENCES categories(id),
    user_id integer REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ,
    UNIQUE (user_id, category_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE budgets
-- +goose StatementEnd
