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
	"github.com/redis/go-redis/v9"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type api struct {
	logger     *zap.Logger
	httpClient *http.Client

	categoryRepo    domain.CategoryRepository
	accountRepo     domain.AccountRepository
	budgetRepo      domain.BudgetRepository
	transactionRepo domain.TransactionRepository
}

func NewAPI(_ context.Context, logger *zap.Logger, _ *redis.Client, pool *pgxpool.Pool) *api {
	categoryRepo := repository.NewPostgresCategory(pool)
	accountRepo := repository.NewPostgresAccount(pool)
	budgetRepo := repository.NewPostgresBudget(pool)
	transctionRepo := repository.NewPostgresTransaction(pool)

	client := &http.Client{}

	return &api{
		logger:     logger,
		httpClient: client,

		categoryRepo:    categoryRepo,
		accountRepo:     accountRepo,
		budgetRepo:      budgetRepo,
		transactionRepo: transctionRepo,
	}
}

func (a *api) Server(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Routes(),
	}
}

func (a *api) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://192.168.254.180:5173", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/health", a.HealthRoutes())
		r.Mount("/categories", a.CategoryRoutes())
		r.Mount("/accounts", a.AccountRoutes())
		r.Mount("/budgets", a.BudgetRoutes())
		r.Mount("/transactions", a.TransactionRoutes())
		r.Mount("/protected", a.ProtectedRoutes())
	})

	return r
}
