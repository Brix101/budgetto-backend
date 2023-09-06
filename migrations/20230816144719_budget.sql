-- +goose Up
-- +goose StatementBegin
-- Add migration script here
CREATE TABLE budgets (
    id SERIAL PRIMARY KEY,
    amount DOUBLE PRECISION DEFAULT 0,
    category_id INTEGER NOT NULL REFERENCES categories (id) ON DELETE NO ACTION,
    created_by VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE,
    UNIQUE (created_by, category_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE budgets
-- +goose StatementEnd
