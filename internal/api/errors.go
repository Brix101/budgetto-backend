package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/go-playground/validator"
)


func (a *api) errorResponse(w http.ResponseWriter, _ *http.Request, status int, err error) {
	var errorResponse domain.ErrResponse

	log.Println(reflect.TypeOf(err), err)
	switch typedErr := err.(type) {
	case error:
		var message string
		if status == 500{
			log.Println(typedErr)
			message =  "Something went wrong!"
		}else{
			message = typedErr.Error()
		}

		errorResponse = domain.ErrResponse{
			Message: message,
			Errors:  []domain.ErrField{},
		}
	case validator.ValidationErrors:
		errorFields := []domain.ErrField{}
		for _, e := range typedErr {
			errorFields = append(errorFields, domain.ErrField{
				Field:   e.Field(),
				Message: GetValidationErrorMessage(e),
			})
		}

		errorResponse = domain.ErrResponse{
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