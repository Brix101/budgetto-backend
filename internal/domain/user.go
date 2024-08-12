package domain

import (
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Base
	Bio      *string `json:"bio,omitempty"`
	Image    *string `json:"image,omitempty"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"-"`
}

func (u *User) NormalizedName() string {
	return strings.ToLower(u.Name)
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// UserRepository represents the user's repository contract
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	// GetAll(ctx context.Context) ([]User, error)

	// CreateOrUpdate(ctx context.Context, usr *User) error
	Update(ctx context.Context, usr *User) (*User, error)
	Create(ctx context.Context, usr *User) (*User, error)
	Delete(ctx context.Context, id int64) error
}
