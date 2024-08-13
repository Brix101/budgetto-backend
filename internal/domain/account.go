package domain

import "context"

type Account struct {
	Base
	Name      string  `json:"name"`
	Note      string  `json:"note,omitempty"`
	CreatedBy uint    `json:"created_by"`
	Balance   float64 `json:"balance"`
}

// AccountRepository represents the account's repository contract
type AccountRepository interface {
	GetByID(ctx context.Context, id uint) (Account, error)
	GetByUserSUB(ctx context.Context, sub string) ([]Account, error)
	// GetAll(ctx context.Context) ([]Account, error)
	//
	// CreateOrUpdate(ctx context.Context, acc *Account) error
	Create(ctx context.Context, acc *Account) (*Account, error)
	Update(ctx context.Context, acc *Account) (*Account, error)
	Delete(ctx context.Context, id int64) error
}
