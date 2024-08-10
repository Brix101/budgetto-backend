package cmd

import (
	"context"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/Brix101/budgetto-backend/internal/api"
	"github.com/Brix101/budgetto-backend/internal/util"
)

func APICmd(ctx context.Context) *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "api",
		Args:  cobra.ExactArgs(0),
		Short: "Runs the RESTful API.",
		RunE: func(_ *cobra.Command, args []string) error {
			port = 5000
			if os.Getenv("PORT") != "" {
				port, _ = strconv.Atoi(os.Getenv("PORT"))
			}

			logger := util.NewLogger("api")
			defer func() { _ = logger.Sync() }()

			db, err := util.NewDatabasePool(ctx, 16)
			if err != nil {
				return err
			}
			defer db.Close()

			redis, err := util.NewRedisQueueClient(ctx, 16)
			if err != nil {
				return err
			}
			defer redis.Close()

			api := api.NewAPI(ctx, logger, redis, db)
			srv := api.Server(port)

			go func() { _ = srv.ListenAndServe() }()

			logger.Info("🚀🚀🚀 Server at port: ", zap.Int("port", port))
			<-ctx.Done()

			_ = srv.Shutdown(ctx)

			return nil
		},
	}

	return cmd
}
