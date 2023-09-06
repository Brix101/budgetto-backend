package domain

import "context"

type Account struct {
	Base

	// account fields
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	Note      string  `json:"note,omitempty"`
	CreatedBy string    `json:"-"`
}

// AccountRepository represents the account's repository contract
type AccountRepository interface {
	GetByID(ctx context.Context, id int64) (Account, error)
	GetByUserSUB(ctx context.Context, sub string) ([]Account, error)
	// GetAll(ctx context.Context) ([]Account, error)
	//
	// CreateOrUpdate(ctx context.Context, acc *Account) error
	Create(ctx context.Context, acc *Account) (*Account, error)
	Update(ctx context.Context, acc *Account) (*Account, error)
	Delete(ctx context.Context, id int64) error
}
