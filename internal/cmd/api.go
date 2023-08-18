package cmd

import (
	"context"

	"github.com/Brix101/budgetto-backend/config"
	"github.com/Brix101/budgetto-backend/internal/api"
	"github.com/Brix101/budgetto-backend/internal/util"
	"github.com/spf13/cobra"
)

func APICmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Args:  cobra.ExactArgs(0),
		Short: "Runs the RESTful API.",
		RunE: func(_ *cobra.Command, _ []string) error {
			env := config.GetConfig()

			db, err := util.NewDatabasePool(ctx, env.DATABASE_URL, 16)
			if err != nil {
				return err
			}
			defer db.Close()

			api := api.NewAPI(db)
			srv := api.Server(env.PORT)

			go func() { _ = srv.ListenAndServe() }()
			<-ctx.Done()

			_ = srv.Shutdown(ctx)

			return nil
		},
	}
	return cmd
}
