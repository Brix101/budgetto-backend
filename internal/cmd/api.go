package cmd

import (
	"context"
	"os"
	"strconv"
	 "log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/Brix101/budgetto-backend/internal/api"
	"github.com/Brix101/budgetto-backend/internal/util"	
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer() func() {
    exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
    if err != nil {
        log.Fatal(err)
    }

    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String("example-service"),
        )),
    )

    otel.SetTracerProvider(tp)

    return func() {
        if err := tp.Shutdown(context.Background()); err != nil {
            log.Fatal(err)
        }
    }
}

func APICmd(ctx context.Context) *cobra.Command {
	var port int
	    cleanup := initTracer()
    		defer cleanup()

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

			if err := util.NewSeeder(ctx, logger, db).CategorySeed(); err != nil {
				return err
			}

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
