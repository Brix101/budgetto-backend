package cmd

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func initTracer() func() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceNameKey.String("budgetto"),
		attribute.String("environment", "development"),
		semconv.ServiceVersionKey.String("0.0.1"),
	))
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		// sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(batchSpanProcessor),
	)
	otel.SetTracerProvider(tracerProvider)

	return func() {
		err := tracerProvider.Shutdown(ctx)
		if err != nil {
			log.Fatalf("failed to shutdown provider: %v", err)
		}
		cancel()
	}
}

func Execute(ctx context.Context) int {
	_ = godotenv.Load()

	cleanup := initTracer()
	defer cleanup()
	
	profile := false

	rootCmd := &cobra.Command{
		Use:   "budgetto",
		Short: "Budgetto is an amazing budgetting app. This isn't it, but it helps.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !profile {
				return nil
			}

			f, perr := os.Create("cpu.pprof")
			if perr != nil {
				return perr
			}

			_ = pprof.StartCPUProfile(f)
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if !profile {
				return nil
			}

			pprof.StopCPUProfile()

			f, perr := os.Create("mem.pprof")
			if perr != nil {
				return perr
			}
			defer f.Close()

			runtime.GC()
			err := pprof.WriteHeapProfile(f)
			return err
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&profile, "profile", "p", false, "record CPU pprof")

	rootCmd.AddCommand(APICmd(ctx))
	rootCmd.AddCommand(SchedulerCmd(ctx))

	go func() {
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	if err := rootCmd.Execute(); err != nil {
		return 1
	}

	return 0
}
