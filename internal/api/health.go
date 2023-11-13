package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func (a api) HealthRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", a.healthCheckHandler)
	r.Get("/protected", a.protectedCheckHandler)

	return r
}

func (hr api) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "available",
		"port":   os.Getenv("PORT"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}

func (a api) protectedCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "available",
		"port":   os.Getenv("PORT"),
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
