-- +goose Up
-- +goose StatementBegin
-- Add migration script here
CREATE TABLE budgets (
    id SERIAL PRIMARY KEY,
    amount DOUBLE PRECISION DEFAULT 0,
    category_id INTEGER UNIQUE REFERENCES categories(id),
    created_by INTEGER REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE,
    UNIQUE (user_id, category_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE budgets
-- +goose StatementEnd
