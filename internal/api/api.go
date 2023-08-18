package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type api struct {
	httpClient *http.Client
}

func NewAPI() *api {
	// tracer := otel.Tracer("api")

	// accountRepo := repository.NewPostgresAccount(pool)
	// deviceRepo := repository.NewPostgresDevice(pool)
	// subredditRepo := repository.NewPostgresSubreddit(pool)
	// watcherRepo := repository.NewPostgresWatcher(pool)
	// userRepo := repository.NewPostgresUser(pool)
	// liveActivityRepo := repository.NewPostgresLiveActivity(pool)

	client := &http.Client{}

	return &api{
		// logger:     logger,
		// statsd:     statsd,
		// reddit:     reddit,
		// apns:       apns,
		httpClient: client,

		// accountRepo:      accountRepo,
		// deviceRepo:       deviceRepo,
		// subredditRepo:    subredditRepo,
		// watcherRepo:      watcherRepo,
		// userRepo:         userRepo,
		// liveActivityRepo: liveActivityRepo,
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
