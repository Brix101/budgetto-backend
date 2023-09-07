package domain

import (
	"context"

	"go.uber.org/zap"
)

type Category struct {
	Base

	// category fields
	Name      string  `json:"name" validate:"required"`
	Note      string  `json:"note,omitempty"`
	CreatedBy *string `json:"created_by,omitempty"`
}

// CategoryRepository represents the categories repository contract
type CategoryRepository interface {
	GetByID(ctx context.Context, id int64) (Category, error)
	GetByUserSUB(ctx context.Context, sub string) ([]Category, error)
	// GetAll(ctx context.Context) ([]Category, error)

	// CreateOrUpdate(ctx context.Context, cat *Category) error
	Update(ctx context.Context, cat *Category) (*Category, error)
	Create(ctx context.Context, cat *Category) (*Category, error)
	Delete(ctx context.Context, id int64) error
	Seed(ctx context.Context, logger *zap.Logger) error
}
