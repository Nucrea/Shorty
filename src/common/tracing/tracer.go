package tracing

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
)

func NewNoopTracer() trace.Tracer {
	return noop.NewTracerProvider().Tracer("shorty")
}

func NewTracer(otelUrl string) (trace.Tracer, error) {
	tracerExporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpointURL(otelUrl),
	)
	if err != nil {
		return nil, err
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("shorty"),
		),
	)

	tracerProvider := traceSdk.NewTracerProvider(
		traceSdk.WithSampler(traceSdk.AlwaysSample()),
		traceSdk.WithResource(r),
		traceSdk.WithBatcher(
			tracerExporter,
			traceSdk.WithMaxQueueSize(8192),
			traceSdk.WithMaxExportBatchSize(2048),
		),
	)

	return tracerProvider.Tracer("shorty"), nil
}
