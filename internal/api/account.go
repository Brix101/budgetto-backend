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

type AccountCtx struct{}

func (a api) AccountRoutes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.Auth)

	r.Get("/", a.accountListHandler)
	r.Post("/", a.accountCreateHandler)

	r.Route("/{id}", func(r chi.Router) {
		r.Use(a.AccountCtx)

		r.Get("/", a.accountGetHandler)
		r.Put("/", a.accountUpdateHandler)
		r.Delete("/", a.accountDeleteHandler)
	})

	return r
}

func (a api) AccountCtx(next http.Handler) http.Handler {
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

		item, err := a.accountRepo.GetByID(ctx, uint(id))
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

type createAccountRequest struct {
	Name    string  `json:"name" validate:"required"`
	Balance float64 `json:"balance" validate:"gte=0"`
	Note    string  `json:"note,omitempty"`
}

func (a api) accountListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value(middlewares.AuthCtx{}).(*domain.UserClaims)
	sub, err := user.GetSubject()
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	accs, err := a.accountRepo.GetByUserSUB(ctx, sub)
	if err != nil {
		a.logger.Error("failed to fetch accounts from database", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(accs)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) accountCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	sub, err := util.GetSub(ctx)
	if err != nil {
		a.errorResponse(w, r, 401, err)
		return
	}

	reqBody := createAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		a.logger.Error("failed to parse request json", zap.Error(err))
		a.errorResponse(w, r, 422, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		a.errorResponse(w, r, 400, err)
		return
	}

	newAcc := domain.Account{
		Name:      reqBody.Name,
		Balance:   reqBody.Balance,
		Note:      reqBody.Note,
		CreatedBy: sub,
	}

	acc, err := a.accountRepo.Create(ctx, &newAcc)
	if err != nil {
		a.logger.Error("failed to create account", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(acc)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) accountGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	item, ok := ctx.Value(AccountCtx{}).(domain.Account)
	if !ok {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	resJSON, err := json.Marshal(item)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) accountUpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	item, ok := ctx.Value(AccountCtx{}).(domain.Account)
	if !ok {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		a.logger.Error("failed to parse request json", zap.Error(err))
		a.errorResponse(w, r, 422, err)
		return
	}

	acc, err := a.accountRepo.Update(ctx, &item)
	if err != nil {
		a.logger.Error("failed to update account", zap.Error(err))
		a.errorResponse(w, r, 500, err)
	}

	resJSON, err := json.Marshal(acc)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) accountDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	item, ok := ctx.Value(AccountCtx{}).(domain.Account)
	if !ok {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	if err := a.accountRepo.Delete(ctx, int64(item.ID)); err != nil {
		a.logger.Error("failed to delete account", zap.Error(err))
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
