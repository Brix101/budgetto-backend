package domain

import "context"

type Category struct {
	Base

	// category fields
	Name      string `json:"name"`
	Note      string `json:"note,omitempty"`
	CreatedBy *uint  `json:"-"`
}

// CategoryRepository represents the categories repository contract
type CategoryRepository interface {
	GetByID(ctx context.Context, id int64) (Category, error)
	GetByUserID(ctx context.Context, user_id int64) ([]Category, error)
	// GetAll(ctx context.Context) ([]Category, error)

	// CreateOrUpdate(ctx context.Context, cat *Category) error
	Update(ctx context.Context, cat *Category) error
	Create(ctx context.Context, cat *Category) error
	Delete(ctx context.Context, id int64) error
}
