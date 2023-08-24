package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)



func (a api) AuthRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/sign-in", a.signInHandler)
	r.Post("/sign-up", a.signUpHandler)

	return r
}

type signInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"` // Minimum length: 6
}

type signUpRequest struct {	
	Name    string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"` // Minimum length: 6
}


func (a api) signInHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var reqBody signInRequest 
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		a.errorResponse(w, r, 422, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqBody)

	if err != nil {
		a.errorResponse(w, r, 400, err)
		return
	}

	usr, err := a.userRepo.GetByEmail(ctx, reqBody.Email)
	if err != nil {		
		a.errorResponse(w, r, 403, domain.ErrInvalidCredentials)
		return
	}

	if validatePass := usr.CheckPassword(reqBody.Password); !validatePass {
		a.errorResponse(w, r, 403, domain.ErrInvalidCredentials)
		return
	}

	usrJSON, err := json.Marshal(usr)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usrJSON)
}

func (a api) signUpHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	reqBody := signUpRequest{} 
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		a.errorResponse(w, r, 422, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqBody)

	if err != nil {
		a.errorResponse(w, r, 400, err)
		return
	}

	newUsr := domain.User{
		Name:reqBody.Name,
		Email:reqBody.Email,
		Password:reqBody.Password,
	}

	if err := newUsr.HashPassword(); err != nil {
		a.errorResponse(w, r, 400, err)
		return
	}

	usr, err := a.userRepo.Create(ctx, &newUsr)
	if err != nil {		
		a.errorResponse(w, r, 500, err)
		return
	}


	usrJSON, err := json.Marshal(usr)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usrJSON)
}