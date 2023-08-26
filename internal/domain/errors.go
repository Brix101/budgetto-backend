package domain

import "errors"

var (
	// ErrNotFound will be returned if the requested item is not found
	ErrNotFound  = errors.New("Requested item was not found.")
	ErrForbidden = errors.New("You don't have permission to access the requested resource.")
	// ErrConflict will be returned if the item being persisted already exists
	ErrConflict = errors.New("Item already exists.")

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
