package domain

import (
	"context"
	"database/sql/driver"
)

type Operation string

const (
	Expense  Operation = "Expense"
	Income   Operation = "Income"
	Transfer Operation = "Transfer"
	Refund   Operation = "Refund"
)

func (st *Operation) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		*st = Operation(b)
	}
	return nil
}

func (st Operation) Value() (driver.Value, error) {
	return string(st), nil
}

type Transaction struct {
	Base

	// transaction fields
	Amount     float64 `json:"amount"`
	Note       string  `json:"note,omitempty"`
	Operation  string  `json:"operation"`
	AccountID  uint    `json:"-"`
	CategoryID uint    `json:"-"`
	CreatedBy  uint    `json:"-"`

	Account  Account  `json:"account,omitempty"`
	Category Category `json:"category,omitempty"`
}

// TransactionRepository represents the transactions repository contract
type TransactionRepository interface {
	GetByID(ctx context.Context, id int64) (Transaction, error)
	GetByUserID(ctx context.Context, user_id int64) ([]Transaction, error)
	// GetAll(ctx context.Context) ([]Transaction, error)

	// CreateOrUpdate(ctx context.Context, tra *Transaction) error
	Update(ctx context.Context, tra *Transaction) (*Transaction, error)
	Create(ctx context.Context, tra *Transaction) (*Transaction, error)
	Delete(ctx context.Context, id int64) error
}
