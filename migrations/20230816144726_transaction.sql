-- +goose Up
-- +goose StatementBegin
DROP TYPE IF EXISTS transaction_type;
CREATE TYPE transaction_type AS ENUM ('Expense', 'Income', 'Transfer', 'Refund');
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    amount DOUBLE PRECISION DEFAULT 0,
    note TEXT,
    transaction_type transaction_type NOT NULL DEFAULT 'Expense',
    account_id integer REFERENCES accounts(id) ON DELETE CASCADE,
    category_id integer REFERENCES categories(id) ON DELETE CASCADE,
    user_id integer REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp,
    deleted_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions
-- +goose StatementEnd
