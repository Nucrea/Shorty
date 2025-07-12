package logging

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/noop"
	sdkLog "go.opentelemetry.io/otel/sdk/log"
)

func newOtelNoop() log.Logger {
	return noop.Logger{}
}

func newOtel(otelUrl string) log.Logger {
	exp, err := otlploghttp.New(
		context.Background(),
		otlploghttp.WithEndpointURL(otelUrl+"/v1/logs"),
	)
	if err != nil {
		panic(err)
	}

	provider := sdkLog.NewLoggerProvider(sdkLog.WithProcessor(sdkLog.NewBatchProcessor(exp)))
	return provider.Logger("shorty")
}
