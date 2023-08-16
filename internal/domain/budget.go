package domain

import "context"

type Budget struct {
	Base

	// budget fields
	Amount     float64 `json:"amount"`
	CategoryID uint    `json:"-"`
	CreatedBy  uint    `json:"-"`
}

// BudgetRepository represents the budget's repository contract
type BudgetRepository interface {
	GetByID(ctx context.Context, id int64) (Budget, error)
	GetByUser(ctx context.Context, id int64) ([]Budget, error)
	GetAll(ctx context.Context) ([]Budget, error)

	CreateOrUpdate(ctx context.Context, acc *Budget) error
	Update(ctx context.Context, acc *Budget) error
	Create(ctx context.Context, acc *Budget) error
	Delete(ctx context.Context, id int64) error
}
