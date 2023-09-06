package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Brix101/budgetto-backend/config"
	"github.com/Brix101/budgetto-backend/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

func (hr api) HealthRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", hr.healthCheckHandler)
	return r
}

func (a api) ProtectedRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares.Auth0Middleware)
	r.Get("/", a.protectedCheckHandler)

	return r
}

func (hr api) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	env := config.GetConfig()
	data := map[string]string{
		"status": "available",
		"port":   env.PORT,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}

func (a api) protectedCheckHandler(w http.ResponseWriter, r *http.Request) {
	env := config.GetConfig()
	data := map[string]string{
		"status": "available",
		"port":   env.PORT,
	}
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()


	claims, err := a.authClaims(ctx)
	if err != nil {
		a.errorResponse(w, r, 403, err)
		return
	}

	fmt.Println(claims.Sub)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}
