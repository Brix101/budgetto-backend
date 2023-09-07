package domain

import (
	"context"
)

type Transaction struct {
	Base

	// transaction fields
	Amount     float64 `json:"amount"`
	Note       string  `json:"note,omitempty"`
	Operation  string  `json:"operation"`
	AccountID  uint    `json:"-"`
	CategoryID uint    `json:"-"`
	CreatedBy  string  `json:"created_by"`

	Account  Account  `json:"account,omitempty"`
	Category Category `json:"category,omitempty"`
}

// TransactionRepository represents the transactions repository contract
type TransactionRepository interface {
	GetByID(ctx context.Context, id int64) (Transaction, error)
	GetByUserSUB(ctx context.Context, sub string) ([]Transaction, error)
	GetOperationType(ctx context.Context) ([]string, error)
	// GetAll(ctx context.Context) ([]Transaction, error)

	// CreateOrUpdate(ctx context.Context, tra *Transaction) error
	Update(ctx context.Context, tra *Transaction) (*Transaction, error)
	Create(ctx context.Context, tra *Transaction) (*Transaction, error)
	Delete(ctx context.Context, id int64) error
}
