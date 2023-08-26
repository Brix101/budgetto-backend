package api

import (
	"fmt"
	"net/http"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/Brix101/budgetto-backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type api struct {
	httpClient *http.Client

	categoryRepo domain.CategoryRepository
	userRepo     domain.UserRepository
	accountRepo  domain.AccountRepository
}

func NewAPI(pool *pgxpool.Pool) *api {
	categoryRepo := repository.NewPostgresCategory(pool)
	userRepo := repository.NewPostgresUser(pool)
	accountRepo := repository.NewPostgresAccount(pool)

	client := &http.Client{}

	return &api{
		httpClient: client,

		categoryRepo: categoryRepo,
		userRepo:     userRepo,
		accountRepo:  accountRepo,
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
	r.Use(middleware.Logger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/health", a.HealthRoutes())
		r.Mount("/categories", a.CategoryRoutes())
		r.Mount("/accounts", a.AccountRoutes())
		r.Mount("/auth", a.AuthRoutes())
		r.Mount("/users", a.UserRoutes())
	})

	return r
}
