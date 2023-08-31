package domain

import "context"

type Budget struct {
	Base

	// budget fields
	Amount     float64 `json:"amount"`
	CategoryID uint    `json:"-"`
	CreatedBy  uint    `json:"-"`

	Category Category `json:"category,omitempty"`
}

// BudgetRepository represents the budget's repository contract
type BudgetRepository interface {
	GetByID(ctx context.Context, id int64) (Budget, error)
	GetByUserID(ctx context.Context, user_id int64) ([]Budget, error)
	// GetAll(ctx context.Context) ([]Budget, error)

	// CreateOrUpdate(ctx context.Context, bud *Budget) error
	Update(ctx context.Context, bud *Budget) (*Budget, error)
	Create(ctx context.Context, bud *Budget) (*Budget, error)
	Delete(ctx context.Context, id int64) error
}
