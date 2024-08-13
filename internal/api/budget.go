package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"go.uber.org/zap"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/Brix101/budgetto-backend/internal/middlewares"
	"github.com/Brix101/budgetto-backend/internal/util"
)

type BudgetCtx struct{}

func (a api) BudgetRoutes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.Auth)

	r.Get("/", a.budgetListHandler)
	r.Post("/", a.budgetCreateHandler)

	r.Route("/{id}", func(r chi.Router) {
		r.Use(a.BudgetCtx)

		r.Get("/{id}", a.budgetGetHandler)
		r.Put("/{id}", a.budgetUpdateHandler)
		r.Delete("/{id}", a.budgetDeleteHandler)
	})

	return r
}

func (a api) BudgetCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sub, err := util.GetSub(ctx)
		if err != nil {
			a.errorResponse(w, r, 500, err)
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			a.errorResponse(w, r, 500, err)
			return
		}

		item, err := a.budgetRepo.GetByID(ctx, uint(id))
		if err != nil {
			status := 500
			if err.Error() == domain.ErrNotFound.Error() {
				status = 404
			}
			a.errorResponse(w, r, status, err)
			return
		}

		if item.CreatedBy != sub {
			a.errorResponse(w, r, 403, domain.ErrForbidden)
			return
		}

		ctx = context.WithValue(ctx, AccountCtx{}, item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type createBudgetRequest struct {
	Amount     float64 `json:"amount" validate:"gte=0"`
	CategoryID uint    `json:"category_id" validate:"required"`
}

func (a api) budgetListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value(middlewares.AuthCtx{}).(*domain.UserClaims)
	sub, err := user.GetSubject()
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	buds, err := a.budgetRepo.GetByUserSUB(ctx, sub)
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

func (a api) budgetCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	sub, err := util.GetSub(ctx)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

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
		CreatedBy:  sub,
	}

	newBud, err := a.budgetRepo.Create(ctx, &budReq)
	if err != nil {
		a.logger.Error("failed to create budget", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	bud, err := a.budgetRepo.GetByID(ctx, newBud.ID)
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

func (a api) budgetGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	bud, ok := ctx.Value(BudgetCtx{}).(domain.Budget)
	if !ok {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
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

	bud := ctx.Value(BudgetCtx{}).(domain.Budget)

	upBud, err := a.budgetRepo.Update(ctx, &bud)
	if err != nil {
		a.logger.Error("failed to update budget", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	cat, err := a.categoryRepo.GetByID(ctx, upBud.CategoryID)
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
		return
	}

	upBud.Category = cat

	resJSON, err := json.Marshal(upBud)
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

	bud, ok := ctx.Value(BudgetCtx{}).(domain.Budget)
	if !ok {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	if err := a.budgetRepo.Delete(ctx, bud.ID); err != nil {
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
