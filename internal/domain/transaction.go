package domain

import (
	"context"
	"database/sql/driver"
)

type TransactionType string

const (
	Expense  TransactionType = "Expense"
	Income   TransactionType = "Income"
	Transfer TransactionType = "Transfer"
	Refund   TransactionType = "Refund"
)

func (st *TransactionType) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		*st = TransactionType(b)
	}
	return nil
}

func (st TransactionType) Value() (driver.Value, error) {
	return string(st), nil
}

type Transaction struct {
	Base

	// transaction fields
	Amount          float64         `json:"amount"`
	Note            string          `json:"note,omitempty"`
	TransactionType TransactionType `json:"transaction_type"`
	AccountID       uint            `json:"-"`
	CategoryID      uint            `json:"-"`
	CreatedBy       uint            `json:"-"`
}

// TransactionRepository represents the transactions repository contract
type TransactionRepository interface {
	GetByID(ctx context.Context, id int64) (Transaction, error)
	GetByUser(ctx context.Context, id int64) ([]Transaction, error)
	GetAll(ctx context.Context) ([]Transaction, error)

	CreateOrUpdate(ctx context.Context, acc *Transaction) error
	Update(ctx context.Context, acc *Transaction) error
	Create(ctx context.Context, acc *Transaction) error
	Delete(ctx context.Context, id int64) error
}
