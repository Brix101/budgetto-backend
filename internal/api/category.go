package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/Brix101/budgetto-backend/internal/middlwares"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

func (a api) CategoryRoutes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlwares.JWTMiddleware)

	r.Get("/", a.categoryListHandler)
	r.Post("/", a.categoryCreateHandler)
	r.Get("/{id}", a.categoryGetHandler)
	r.Put("/{id}", a.categoryUpdateHandler)
	r.Delete("/{id}", a.categoryDeleteHandler)

	return r
}

type updateCategoryRequest struct {
	Name string `json:"name"`
	Note string `json:"note"`
}

func (a api) categoryListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	cats, err := a.categoryRepo.GetByUserID(ctx, int64(user.Sub))
	if err != nil {
		a.logger.Error("failed to fetch categories from database", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(cats)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) categoryGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	cat, err := a.categoryRepo.GetByID(ctx, int64(id))
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
		return
	}

	if cat.ID != uint(user.Sub) {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
		return
	}

	resJSON, err := json.Marshal(cat)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) categoryCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	userId := uint(user.Sub)
	newCat := domain.Category{
		CreatedBy: &userId,
	}

	if err := json.NewDecoder(r.Body).Decode(&newCat); err != nil {
		a.logger.Error("failed to parse request json", zap.Error(err))
		a.errorResponse(w, r, 422, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(newCat); err != nil {
		a.logger.Error("failed to validate create category struct", zap.Error(err))
		a.errorResponse(w, r, 400, err)
		return
	}

	cat, err := a.categoryRepo.Create(ctx, &newCat)
	if err != nil {
		a.logger.Error("failed to create category", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	catJSON, err := json.Marshal(cat)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(catJSON)
}

func (a api) categoryUpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	cat, err := a.categoryRepo.GetByID(ctx, int64(id))
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
		return
	}

	if cat.ID != uint(user.Sub) {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
		return
	}

	var upCat updateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&upCat); err != nil {
		a.logger.Error("failed to parse request json", zap.Error(err))
		a.errorResponse(w, r, 422, err)
		return
	}
	defer r.Body.Close()

	if upCat.Name != "" {
		cat.Name = upCat.Name
	}

	if upCat.Note != "" {
		cat.Note = upCat.Note
	}

	updatedCat, err := a.categoryRepo.Update(ctx, &cat)
	if err != nil {
		a.logger.Error("failed to delete category", zap.Error(err))
		a.errorResponse(w, r, 500, err)
	}

	catJSON, err := json.Marshal(updatedCat)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(catJSON)
}

func (a api) categoryDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	cat, err := a.categoryRepo.GetByID(ctx, int64(id))
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
		return
	}

	if cat.ID != uint(user.Sub) {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
		return
	}
	if err := a.categoryRepo.Delete(ctx, int64(id)); err != nil {
		a.logger.Error("failed to delete category", zap.Error(err))
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
	}

	data := map[string]string{
		"message": "Item deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}
