package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/Brix101/budgetto-backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type api struct {
	logger     *zap.Logger
	httpClient *http.Client

	categoryRepo domain.CategoryRepository
	userRepo     domain.UserRepository
	accountRepo  domain.AccountRepository
	budgetRepo   domain.BudgetRepository
}

func NewAPI(ctx context.Context, logger *zap.Logger, pool *pgxpool.Pool) *api {
	categoryRepo := repository.NewPostgresCategory(pool)
	userRepo := repository.NewPostgresUser(pool)
	accountRepo := repository.NewPostgresAccount(pool)
	budgetRepo := repository.NewPostgresBudget(pool)

	client := &http.Client{}

	categoryRepo.Seed(ctx)

	return &api{
		logger:     logger,
		httpClient: client,

		categoryRepo: categoryRepo,
		userRepo:     userRepo,
		accountRepo:  accountRepo,
		budgetRepo:   budgetRepo,
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
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/health", a.HealthRoutes())
		r.Mount("/categories", a.CategoryRoutes())
		r.Mount("/accounts", a.AccountRoutes())
		r.Mount("/auth", a.AuthRoutes())
		r.Mount("/users", a.UserRoutes())
		r.Mount("/budgets", a.BudgetRoutes())
	})

	return r
}
