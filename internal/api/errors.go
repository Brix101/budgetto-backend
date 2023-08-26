package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgconn"
)

func (a *api) errorResponse(w http.ResponseWriter, _ *http.Request, status int, err error) {
	var errorResponse domain.ErrResponse
	statusCode := status // Default status code

	log.Println(reflect.TypeOf(err), err)
	switch typedErr := err.(type) {
	case *pgconn.PgError:
		constraintName := typedErr.ConstraintName
		parts := strings.Split(constraintName, "_")
		fieldName := "field" // Default field name if constraint name is not in the expected format
		if len(parts) > 0 {
			fieldName = parts[len(parts)-1]
			if fieldName == "key" {
				fieldName = parts[1]
			}
		}

		switch typedErr.Code {
		case "23505":
			statusCode = http.StatusBadRequest
			message := fmt.Sprintf("The %s you entered is already taken.", fieldName)
			errorFields := []domain.ErrField{{
				Field:   fieldName,
				Message: message,
			}}

			errorResponse = domain.ErrResponse{
				Message: "Validation Error",
				Errors:  errorFields,
			}

		// case "23503":
		// 	statusCode = http.StatusBadRequest
		// 	// message := fmt.Sprintf("%s is already taken.", strings.Title(fieldName))
		// 	return statusCode, []ErrorInfo{{Field: fieldName, Message: typedErr.Message}}
		// Add more cases to handle other PostgreSQL error codes if needed
		// ...
		default:
			statusCode = 500
			errorResponse = domain.ErrResponse{
				Message: "Something went wrong!",
				Errors:  []domain.ErrField{},
			}
		}

	case validator.ValidationErrors:
		fmt.Println("Validation")
		errorFields := []domain.ErrField{}
		for _, e := range typedErr {
			errorFields = append(errorFields, domain.ErrField{
				Field:   strings.ToLower(e.Field()),
				Message: GetValidationErrorMessage(e),
			})
		}

		errorResponse = domain.ErrResponse{
			Message: "Validation Error",
			Errors:  errorFields,
		}

	case error:
		var message string
		if statusCode == 500 {
			message = "Something went wrong!"
		} else {
			message = typedErr.Error()
		}
		errorResponse = domain.ErrResponse{
			Message: message,
			Errors:  []domain.ErrField{},
		}

	default:
		statusCode = 500
		errorResponse = domain.ErrResponse{
			Message: "Something went wrong!",
			Errors:  []domain.ErrField{},
		}

	}

	errorResJson, _ := json.Marshal(errorResponse)

	w.Header().Set("X-Budgetto-Error", err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(errorResJson)
}

func GetValidationErrorMessage(err validator.FieldError) string {
	fieldName := strings.ToLower(err.Field())

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s field is required.", fieldName)
	case "min":
		return fmt.Sprintf("%s should be at least %s characters long.", fieldName, err.Param())
	case "max":
		return fmt.Sprintf("%s should be at most %s characters long.", fieldName, err.Param())
	case "email":
		return "Enter a valid email address."
	case "gte":
		return fmt.Sprintf("%s should be greater than %s.", fieldName, err.Param())
	// Add more cases for other validation tags as needed.
	default:
		return "Invalid input."
	}
}
