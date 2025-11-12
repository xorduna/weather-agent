package observability

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"log"
	"time"
)

var (
	meter                 metric.Meter
	ResponseTimeHistogram metric.Float64Histogram
	RequestCounter        metric.Int64Counter
)

func InitMetrics(ctx context.Context) func() {
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create OTLP exporter: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
	)
	otel.SetMeterProvider(provider)

	meter = otel.Meter("chat-server")

	RequestCounter, err = meter.Int64Counter("requests.total",
		metric.WithDescription("Total requests processed"))
	if err != nil {
		panic(err)
	}

	ResponseTimeHistogram, err = meter.Float64Histogram("requests.response_time",
		metric.WithDescription("Response time in milliseconds"),
		metric.WithUnit("ms"))
	if err != nil {
		panic(err)
	}

	return func() {
		if err := provider.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown meter provider: %v", err)
		}
	}
}

func RecordRequest(ctx context.Context, method string, start time.Time, errPtr *error) {
	duration := time.Since(start)
	success := errPtr == nil || *errPtr == nil

	attrs := []attribute.KeyValue{
		attribute.String("method", method),
		attribute.Bool("success", success),
	}
	RequestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	ResponseTimeHistogram.Record(ctx, float64(duration.Milliseconds()), metric.WithAttributes(attrs...))
}
