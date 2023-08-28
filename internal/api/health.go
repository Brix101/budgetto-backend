package api

import (
	"encoding/json"
	"net/http"

	"github.com/Brix101/budgetto-backend/config"
	"github.com/go-chi/chi/v5"
)

func (hr api) HealthRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", hr.healthCheckHandler)
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
