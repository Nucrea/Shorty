package tracing

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/trace"

	traceSdk "go.opentelemetry.io/otel/sdk/trace"
)

func NewTracer(otelUrl string) (trace.Tracer, error) {
	tracerExporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpointURL(otelUrl),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := traceSdk.NewTracerProvider(
		traceSdk.WithSampler(traceSdk.AlwaysSample()),
		traceSdk.WithBatcher(
			tracerExporter,
			traceSdk.WithMaxQueueSize(8192),
			traceSdk.WithMaxExportBatchSize(2048),
		),
	)

	return tracerProvider.Tracer("shorty"), nil
}
