package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"go.uber.org/zap"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/Brix101/budgetto-backend/internal/middlewares"
)

type CatCtx struct{}

func (a api) CategoryRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares.AuthMiddleware)

	r.Get("/", a.categoryListHandler)
	r.Post("/", a.categoryCreateHandler)

	r.Route("/{id}", func(r chi.Router) {
		r.Use(a.CategoryCtx)

		r.Get("/", a.categoryGetHandler)
		r.Put("/", a.categoryUpdateHandler)
		r.Delete("/", a.categoryDeleteHandler)
	})

	return r
}

func (a api) CategoryCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		fmt.Println(r.Method)
		user := ctx.Value(middlewares.UserCtxKey{}).(*domain.UserClaims)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			a.errorResponse(w, r, 500, err)
			return
		}

		item, err := a.categoryRepo.GetByID(ctx, int64(id))
		if err != nil {
			status := 500
			if err.Error() == domain.ErrNotFound.Error() {
				status = 404
			}
			a.errorResponse(w, r, status, err)
			return
		}

		if item.CreatedBy != nil && *item.CreatedBy != uint(user.Sub) {
			a.errorResponse(w, r, 403, domain.ErrForbidden)
			return
		}

		ctx = context.WithValue(ctx, CatCtx{}, item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type createCategoryRequest struct {
	Name string `json:"name" validate:"required"`
	Note string `json:"note,omitempty"`
}

func (a api) categoryListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value(middlewares.UserCtxKey{}).(*domain.UserClaims)

	cats, err := a.categoryRepo.GetByUserSUB(ctx, int64(user.Sub))
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

func (a api) categoryCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value(middlewares.UserCtxKey{}).(*domain.UserClaims)

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

	userId := uint(user.Sub)
	newCat := domain.Category{
		Name:      reqBody.Name,
		Note:      reqBody.Note,
		CreatedBy: &userId,
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

func (a api) categoryGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	item, ok := ctx.Value(CatCtx{}).(domain.Category)
	if !ok {

		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	res, err := json.Marshal(item)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (a api) categoryUpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value(middlewares.UserCtxKey{}).(*domain.UserClaims)
	item, ok := ctx.Value(CatCtx{}).(domain.Category)
	if !ok {

		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	if item.CreatedBy == nil || *item.CreatedBy != uint(user.Sub) {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		a.logger.Error("failed to parse request json", zap.Error(err))
		a.errorResponse(w, r, 422, err)
		return
	}
	defer r.Body.Close()

	updatedCat, err := a.categoryRepo.Update(ctx, &item)
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
	user := r.Context().Value(middlewares.UserCtxKey{}).(*domain.UserClaims)
	item, ok := ctx.Value(CatCtx{}).(domain.Category)
	if !ok {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	if item.CreatedBy == nil || *item.CreatedBy != uint(user.Sub) {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
		return
	}
	if err := a.categoryRepo.Delete(ctx, int64(item.ID)); err != nil {
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
