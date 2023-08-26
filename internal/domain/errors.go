package domain

import "errors"

var (
	ErrNotFound           = errors.New("Requested item was not found.")
	ErrForbidden          = errors.New("You don't have permission to access the requested resource.")
	ErrInvalidCredentials = errors.New("Invalid credentials. Please try again.")
)

type ErrResponse struct {
	Message string     `json:"message"`
	Errors  []ErrField `json:"errors,omitempty"`
}

type ErrField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
