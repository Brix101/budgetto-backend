package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

func (cr api) CategoryRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", cr.categoryListHandler)
	r.Post("/", cr.categoryCreateHandler)
	return r
}

func (cr api) categoryListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	categories, err := cr.categoryRepo.GetByUserID(ctx, 5)
	if err != nil {
		log.Fatal(err)
	}

	catsJson, err := json.Marshal(categories)
	if err != nil {
		cr.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(catsJson)
}

func (cr api) categoryCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var cat domain.Category
	userID := uint(1)
	cat.CreatedBy = &userID
	cat.Note = ""

	err := json.NewDecoder(r.Body).Decode(&cat)
	if err != nil {
		cr.errorResponse(w, r, 422, err)
		return
	}
	// Validate the user struct
	validate := validator.New()
	err = validate.Struct(cat)

	if err != nil {
		cr.errorResponse(w, r, 400, err)
		return
	}

	category, err := cr.categoryRepo.Create(ctx, &cat)
	if err != nil {
		log.Fatal(err)
	}

	catJson, err := json.Marshal(category)
	if err != nil {
		cr.errorResponse(w, r, 500, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(catJson)
}
