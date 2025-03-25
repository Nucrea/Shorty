package metrics

import "go.opentelemetry.io/otel/metric/noop"

func NewNoop() Meter {
	return otelMetrics{noop.Meter{}}
}
