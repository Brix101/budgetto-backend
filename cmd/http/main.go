package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Brix101/budgetto-backend/config"
	"github.com/Brix101/budgetto-backend/internal/api"
	"github.com/Brix101/budgetto-backend/internal/util"
	"go.uber.org/zap"
)

func main() {
	env := config.GetConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	logger := util.NewLogger("api")
	defer func() { _ = logger.Sync() }()

	db, err := util.NewDatabasePool(ctx, env.DATABASE_URL, 16)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	api := api.NewAPI(ctx, logger, db)
	server := api.Server(env.PORT)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
		defer shutdownCancel() // Call the cancel function when the shutdown function finishes

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		cancel()
	}()

	logger.Info("🚀🚀🚀 Server at port: " + env.PORT)
	// Run the server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Server failed to start: ", zap.Error(err))
	}

	// Wait for server context to be stopped
	<-ctx.Done()
}
