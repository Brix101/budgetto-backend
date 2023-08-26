package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/Brix101/budgetto-backend/internal/middlwares"
	"github.com/go-chi/chi/v5"
)

func (a api) AccountRoutes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlwares.JWTMiddleware)

	r.Get("/", a.accountListHandler)
	r.Get("/{id}", a.accountGetHandler)

	return r
}

func (a api) accountListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value("user").(*domain.UserClaims)

	accs, err := a.accountRepo.GetByUserID(ctx, int64(user.Sub))
	if err != nil {
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

	resJSON, err := json.Marshal(acc)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}
