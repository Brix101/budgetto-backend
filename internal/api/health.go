package api

import (
	"encoding/json"
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}
