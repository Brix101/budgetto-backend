package api

import (
	"fmt"
	"net/http"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/Brix101/budgetto-backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type api struct {
	httpClient *http.Client

	categoryRepo domain.CategoryRepository
}

func NewAPI(pool *pgxpool.Pool) *api {
	categoryRepo := repository.NewPostgresCategory(pool)

	client := &http.Client{}

	return &api{
		httpClient: client,

		categoryRepo: categoryRepo,
	}
}

func (a *api) Server(port string) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: a.Routes(),
	}
}

func (a *api) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/health", a.HealthRoutes())
	})

	return r
}
