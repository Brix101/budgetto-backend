package domain

import "context"

type Budget struct {
	Base
	Category   Category `json:"category,omitempty"`
	CreatedBy  uint     `json:"created_by"`
	Amount     float64  `json:"amount"`
	CategoryID uint     `json:"-"`
}

// BudgetRepository represents the budget's repository contract
type BudgetRepository interface {
	GetByID(ctx context.Context, id uint) (Budget, error)
	GetByUserSUB(ctx context.Context, sub string) ([]Budget, error)
	// GetAll(ctx context.Context) ([]Budget, error)

	// CreateOrUpdate(ctx context.Context, bud *Budget) error
	Update(ctx context.Context, bud *Budget) (*Budget, error)
	Create(ctx context.Context, bud *Budget) (*Budget, error)
	Delete(ctx context.Context, id uint) error
}
