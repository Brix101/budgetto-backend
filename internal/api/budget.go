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

func (a api) BudgetRoutes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.AuthMiddleware)

	r.Get("/", a.budgetListHandler)
	r.Post("/", a.budgetCreateHandler)
	r.Get("/{id}", a.budgetGetHandler)
	r.Put("/{id}", a.budgetUpdateHandler)
	r.Delete("/{id}", a.budgetDeleteHandler)

	return r
}

type createBudgetRequest struct {
	Amount     float64 `json:"amount" validate:"gte=0"`
	CategoryID uint    `json:"category_id" validate:"required"`
}

func (a api) budgetListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	buds, err := a.budgetRepo.GetByUserID(ctx, int64(user.Sub))
	if err != nil {
		a.logger.Error("failed to fetch budgets from database", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(buds)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) budgetGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	bud, err := a.budgetRepo.GetByID(ctx, int64(id))
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		if status == 500 {
			a.logger.Error("failed to fetch from database", zap.Error(err))
		}
		a.errorResponse(w, r, status, err)
		return
	}

	if bud.CreatedBy != uint(user.Sub) {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
		return
	}

	resJSON, err := json.Marshal(bud)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) budgetCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	userId := uint(user.Sub)
	reqBody := createBudgetRequest{}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		a.logger.Error("failed to parse request json", zap.Error(err))
		a.errorResponse(w, r, 422, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		a.logger.Error("failed to validate create budget struct", zap.Error(err))
		a.errorResponse(w, r, 400, err)
		return
	}

	budReq := domain.Budget{
		Amount:     reqBody.Amount,
		CategoryID: reqBody.CategoryID,
		CreatedBy:  userId,
	}

	newBud, err := a.budgetRepo.Create(ctx, &budReq)
	if err != nil {
		a.logger.Error("failed to create budget", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	bud, err := a.budgetRepo.GetByID(ctx, int64(newBud.ID))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(bud)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) budgetUpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	bud, err := a.budgetRepo.GetByID(ctx, int64(id))
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
		return
	}

	if bud.CreatedBy != uint(user.Sub) {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&bud); err != nil {
		a.logger.Error("failed to parse request json", zap.Error(err))
		a.errorResponse(w, r, 422, err)
		return
	}
	defer r.Body.Close()

	updatedBud, err := a.budgetRepo.Update(ctx, &bud)
	if err != nil {
		a.logger.Error("failed to update budget", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	cat, err := a.categoryRepo.GetByID(ctx, int64(updatedBud.CategoryID))
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
		return
	}

	updatedBud.Category = cat

	resJSON, err := json.Marshal(updatedBud)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) budgetDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	bud, err := a.budgetRepo.GetByID(ctx, int64(id))
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
		return
	}

	if bud.CreatedBy != uint(user.Sub) {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
		return
	}

	if err := a.budgetRepo.Delete(ctx, int64(id)); err != nil {
		a.logger.Error("failed to delete budget", zap.Error(err))
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
