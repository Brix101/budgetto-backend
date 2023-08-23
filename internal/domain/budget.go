package domain

import "context"

type Budget struct {
	Base

	// budget fields
	Amount     float64 `json:"amount"`
	CategoryID uint    `json:"category_id"`
	CreatedBy  uint    `json:"created_by"`
}

// BudgetRepository represents the budget's repository contract
type BudgetRepository interface {
	GetByID(ctx context.Context, id int64) (Budget, error)
	GetByUser(ctx context.Context, id int64) ([]Budget, error)
	GetAll(ctx context.Context) ([]Budget, error)

	CreateOrUpdate(ctx context.Context, bud *Budget) error
	Update(ctx context.Context, bud *Budget) error
	Create(ctx context.Context, bud *Budget) error
	Delete(ctx context.Context, id int64) error
}
