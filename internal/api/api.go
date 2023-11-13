package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/Brix101/budgetto-backend/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type api struct {
	logger     *zap.Logger
	httpClient *http.Client

	categoryRepo    domain.CategoryRepository
	accountRepo     domain.AccountRepository
	budgetRepo      domain.BudgetRepository
	transactionRepo domain.TransactionRepository
	userRepo        domain.UserRepository
}

func NewAPI(_ context.Context, logger *zap.Logger, _ *redis.Client, pool *pgxpool.Pool) *api {
	categoryRepo := repository.NewPostgresCategory(pool)
	accountRepo := repository.NewPostgresAccount(pool)
	budgetRepo := repository.NewPostgresBudget(pool)
	transctionRepo := repository.NewPostgresTransaction(pool)
	userRepo := repository.NewPostgresUser(pool)

	client := &http.Client{}

	return &api{
		logger:     logger,
		httpClient: client,

		categoryRepo:    categoryRepo,
		accountRepo:     accountRepo,
		budgetRepo:      budgetRepo,
		transactionRepo: transctionRepo,
		userRepo:        userRepo,
	}
}

func (a *api) Server(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Routes(),
	}
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v1

func (a *api) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://192.168.254.180:5173", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:1323/swagger/doc.json"), // The url pointing to API definition
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/health", a.HealthRoutes())
		r.Mount("/categories", a.CategoryRoutes())
		r.Mount("/accounts", a.AccountRoutes())
		r.Mount("/budgets", a.BudgetRoutes())
		r.Mount("/transactions", a.TransactionRoutes())
		r.Mount("/auth", a.AuthRoutes())
		r.Mount("/users", a.UserRoutes())
	})

	return r
}
