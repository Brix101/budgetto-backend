package domain

import (
	"context"
)

type Transaction struct {
	Base
	Category   Category `json:"category,omitempty"`
	Note       string   `json:"note,omitempty"`
	Operation  string   `json:"operation"`
	CreatedBy  int      `json:"created_by"`
	Account    Account  `json:"account,omitempty"`
	Amount     float64  `json:"amount"`
	AccountID  uint     `json:"-"`
	CategoryID uint     `json:"-"`
}

// TransactionRepository represents the transactions repository contract
type TransactionRepository interface {
	GetByID(ctx context.Context, id int64) (Transaction, error)
	GetByUserSUB(ctx context.Context, sub int64) ([]Transaction, error)
	GetOperationType(ctx context.Context) ([]string, error)
	// GetAll(ctx context.Context) ([]Transaction, error)

	// CreateOrUpdate(ctx context.Context, tra *Transaction) error
	Update(ctx context.Context, tra *Transaction) (*Transaction, error)
	Create(ctx context.Context, tra *Transaction) (*Transaction, error)
	Delete(ctx context.Context, id int64) error
}
