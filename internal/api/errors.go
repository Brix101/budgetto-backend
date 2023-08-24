package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
)

type ErrorResponse struct {
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (a *api) errorResponse(w http.ResponseWriter, _ *http.Request, status int, err error) {
	var errorResponse ErrorResponse

	switch typedErr := err.(type) {
	case error:
		var message string
		if status == 500{
			log.Println(typedErr)
			message =  "Something went wrong!"
		}else{
			message = typedErr.Error()
		}

		errorResponse = ErrorResponse{
			Message: message,
			Errors:  []FieldError{},
		}
	case validator.ValidationErrors:
		errorFields := []FieldError{}
		for _, e := range typedErr {
			errorFields = append(errorFields, FieldError{
				Field:   e.Field(),
				Message: GetValidationErrorMessage(e),
			})
		}

		errorResponse = ErrorResponse{
			Message: "Validation Error",
			Errors:  errorFields,
		}
	}

	errorResJson, _ := json.Marshal(errorResponse)

	// w.Header().Set("X-Budgetto-Error", err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(errorResJson)
}

func GetValidationErrorMessage(err validator.FieldError) string {
	fieldName := err.Field()
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s field is required.", fieldName)
	case "min":
		minValue, _ := strconv.Atoi(err.Param())
		return fmt.Sprintf("%s should be at least %d characters long.", fieldName, minValue)
	case "max":
		maxValue, _ := strconv.Atoi(err.Param())
		return fmt.Sprintf("%s should be at most %d characters long.", fieldName, maxValue)
	case "email":
		return "Enter a valid email address."
	// Add more cases for other validation tags as needed.
	default:
		return "Invalid input."
	}
}