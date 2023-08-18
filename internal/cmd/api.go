package cmd

import (
	"context"

	"github.com/Brix101/budgetto-backend/config"
	"github.com/Brix101/budgetto-backend/internal/api"
	"github.com/spf13/cobra"
)

func APICmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Args:  cobra.ExactArgs(0),
		Short: "Runs the RESTful API.",
		RunE: func(_ *cobra.Command, _ []string) error {
			env := config.GetConfig()

			api := api.NewAPI()
			srv := api.Server(env.PORT)

			go func() { _ = srv.ListenAndServe() }()
			<-ctx.Done()

			_ = srv.Shutdown(ctx)

			return nil
		},
	}
	return cmd
}
