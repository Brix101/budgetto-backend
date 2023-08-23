package api

import (
	"errors"
	"net/http"
)

var ErrDuplicateAPNSToken = errors.New("duplicate apns token")

type ErrorResponse struct {
	Message string      `json:"message"`
	Errors  []FieldError  `json:"errors"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (a *api) errorResponse(w http.ResponseWriter, _ *http.Request, status int, err error) {
	w.Header().Set("X-Budgetto-Error", err.Error())
	http.Error(w, err.Error(), status)
}
