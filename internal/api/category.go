package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/Brix101/budgetto-backend/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

func (a api) CategoryRoutes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.Auth0Middleware)

	r.Get("/", a.categoryListHandler)
	r.Post("/", a.categoryCreateHandler)
	r.Get("/{id}", a.categoryGetHandler)
	r.Put("/{id}", a.categoryUpdateHandler)
	r.Delete("/{id}", a.categoryDeleteHandler)

	return r
}

type createCategoryRequest struct {
	Name string `json:"name" validate:"required"`
	Note string `json:"note,omitempty"`
}

func (a api) categoryListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user,err := a.authClaims(ctx)
	if err != nil {
		a.errorResponse(w, r, 403, err)
		return
	}

	cats, err := a.categoryRepo.GetByUserSUB(ctx, user.Sub)
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

	user,err := a.authClaims(ctx)
	if err != nil {
		a.errorResponse(w, r, 403, err)
		return
	}

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

	if *cat.CreatedBy != user.Sub {
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

	user,err := a.authClaims(ctx)
	if err != nil {
		a.errorResponse(w, r, 403, err)
		return
	}

	reqBody := createCategoryRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		a.logger.Error("failed to parse request json", zap.Error(err))
		a.errorResponse(w, r, 422, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		a.logger.Error("failed to validate create category struct", zap.Error(err))
		a.errorResponse(w, r, 400, err)
		return
	}

	newCat := domain.Category{
		Name:      reqBody.Name,
		Note:      reqBody.Note,
		CreatedBy: &user.Sub,
	}

	cat, err := a.categoryRepo.Create(ctx, &newCat)
	if err != nil {
		a.logger.Error("failed to create category", zap.Error(err))
		a.errorResponse(w, r, 500, err)
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

func (a api) categoryUpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user,err := a.authClaims(ctx)
	if err != nil {
		a.errorResponse(w, r, 403, err)
		return
	}

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

	if *cat.CreatedBy != user.Sub {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		a.logger.Error("failed to parse request json", zap.Error(err))
		a.errorResponse(w, r, 422, err)
		return
	}
	defer r.Body.Close()

	updatedCat, err := a.categoryRepo.Update(ctx, &cat)
	if err != nil {
		a.logger.Error("failed to delete category", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(updatedCat)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) categoryDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	user,err := a.authClaims(ctx)
	if err != nil {
		a.errorResponse(w, r, 403, err)
		return
	}

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

	if *cat.CreatedBy != user.Sub {
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
