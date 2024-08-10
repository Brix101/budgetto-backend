package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Brix101/budgetto-backend/internal/util"
)

func SeedCmd(ctx context.Context) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "seed",
		Args:  cobra.ExactArgs(0),
		Short: "Run the seed command to populate the database with initial data.",
		RunE: func(_ *cobra.Command, args []string) error {
			logger := util.NewLogger("api")
			defer func() { _ = logger.Sync() }()

			db, err := util.NewDatabasePool(ctx, 16)
			if err != nil {
				return err
			}
			defer db.Close()

			if err := util.NewSeeder(ctx, logger, db).CategorySeed(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
