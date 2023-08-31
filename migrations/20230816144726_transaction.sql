-- +goose Up
-- +goose StatementBegin
DROP TYPE IF EXISTS operation;
CREATE TYPE operation AS ENUM ('Expense', 'Income', 'Transfer', 'Refund');
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    amount DOUBLE PRECISION DEFAULT 0,
    note TEXT DEFAULT '',
    operation operation NOT NULL DEFAULT 'Expense',
    account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_by INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions
-- +goose StatementEnd
