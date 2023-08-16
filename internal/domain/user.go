package domain

import "strings"

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
