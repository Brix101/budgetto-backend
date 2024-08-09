package util

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func InitTracer() func() {

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
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
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
