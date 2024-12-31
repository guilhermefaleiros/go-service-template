package observability

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"guilhermefaleiros/go-service-template/internal/shared"
	"log/slog"
)

var (
	meterProvider *metric.MeterProvider
)

func InitMeterProvider(cfg *shared.Config) {
	exporter, err := prometheus.New()
	if err != nil {
		slog.Error("failed to create prometheus exporter")
	}

	meterProvider = metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(resource.Default()),
	)

	otel.SetMeterProvider(meterProvider)

	meterProvider.Meter(cfg.App.Name)
}
