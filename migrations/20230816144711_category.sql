-- +goose Up
-- +goose StatementBegin
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    note TEXT,
    created_by INTEGER REFERENCES users (id) ON DELETE CASCADE DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE categories
-- +goose StatementEnd
