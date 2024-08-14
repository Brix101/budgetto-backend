package budgettocmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/Brix101/budgetto-backend/internal/util"
)

func SchedulerCmd(ctx context.Context) *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "scheduler",
		Args:  cobra.ExactArgs(0),
		Short: "Schedules jobs and runs several maintenance tasks periodically.",
		RunE: func(_ *cobra.Command, args []string) error {
			port = 5000
			if os.Getenv("PORT") != "" {
				port, _ = strconv.Atoi(os.Getenv("PORT"))
			}

			logger := util.NewLogger("api")
			defer func() { _ = logger.Sync() }()

			fmt.Println(port, args)
			s := gocron.NewScheduler(time.UTC)
			s.SetMaxConcurrentJobs(8, gocron.WaitMode)

			_, _ = s.Every(5).Seconds().Do(func() { enqueueLiveActivities(ctx, logger) })
			s.StartAsync()

			srv := &http.Server{Addr: ":8080"}
			go func() { _ = srv.ListenAndServe() }()

			<-ctx.Done()

			s.Stop()

			return nil
		},
	}

	return cmd
}

func enqueueLiveActivities(_ context.Context, logger *zap.Logger) {
	logger.Info("Pinging ....")
}
