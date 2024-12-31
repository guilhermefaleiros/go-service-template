package observability

import (
	"context"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"guilhermefaleiros/go-service-template/internal/shared"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer(cfg *shared.Config) *sdktrace.TracerProvider {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(cfg.Observability.Jaeger.Endpoint),
	))
	if err != nil {
		log.Fatalf("Error creating Jaeger exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			"",
			attribute.String("service.name", cfg.App.Name),
			attribute.String("environment", cfg.App.Environment),
			attribute.String("version", "1.0.0"),
		)),
	)

	otel.SetTracerProvider(tp)

	return tp
}

func ShutdownTracerProvider(tp *sdktrace.TracerProvider) {
	if err := tp.Shutdown(context.Background()); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	}
}
