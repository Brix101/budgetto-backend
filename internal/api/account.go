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

func (a api) AccountRoutes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.AuthMiddleware)

	r.Get("/", a.accountListHandler)
	r.Post("/", a.accountCreateHandler)
	r.Get("/{id}", a.accountGetHandler)
	r.Put("/{id}", a.accountUpdateHandler)
	r.Delete("/{id}", a.accountDeleteHandler)

	return r
}

type createAccountRequest struct {
	Name    string  `json:"name" validate:"required"`
	Balance float64 `json:"balance" validate:"gte=0"`
	Note    string  `json:"note,omitempty"`
}

func (a api) accountListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	accs, err := a.accountRepo.GetByUserID(ctx, int64(user.Sub))
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

func (a api) accountGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	acc, err := a.accountRepo.GetByID(ctx, int64(id))
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
		return
	}

	if acc.CreatedBy != uint(user.Sub) {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
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

func (a api) accountCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

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
		CreatedBy: uint(user.Sub),
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

func (a api) accountUpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	acc, err := a.accountRepo.GetByID(ctx, int64(id))
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
		return
	}

	if acc.CreatedBy != uint(user.Sub) {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		a.logger.Error("failed to parse request json", zap.Error(err))
		a.errorResponse(w, r, 422, err)
		return
	}

	updatedCat, err := a.accountRepo.Update(ctx, &acc)
	if err != nil {
		a.logger.Error("failed to update account", zap.Error(err))
		a.errorResponse(w, r, 500, err)
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

func (a api) accountDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	acc, err := a.accountRepo.GetByID(ctx, int64(id))
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
		return
	}

	if acc.CreatedBy != uint(user.Sub) {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
		return
	}

	if err = a.accountRepo.Delete(ctx, int64(acc.ID)); err != nil {
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
