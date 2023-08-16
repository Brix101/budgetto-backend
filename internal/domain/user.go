package domain

import (
	"context"
	"strings"
)

type User struct {
	Base

	// user fields
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"-"`
	Bio      *string `json:"bio,omitempty"`
	Image    *string `json:"image,omitempty"`
}

func (u *User) NormalizedName() string {
	return strings.ToLower(u.Name)
}

func (u User) CheckPassword(password string) bool {
	if u.Password == password {
		return true
	}
	return false
}

// UserRepository represents the user's repository contract
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetAll(ctx context.Context) ([]User, error)

	CreateOrUpdate(ctx context.Context, acc *User) error
	Update(ctx context.Context, acc *User) error
	Create(ctx context.Context, acc *User) error
	Delete(ctx context.Context, id int64) error
}
