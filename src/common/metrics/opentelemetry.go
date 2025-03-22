package metrics

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
)

var _ Meter = (*otelMetrics)(nil)

func NewOtel(prefix string, otelUrl string) *otelMetrics {
	exporter, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpointURL(otelUrl),
	)
	if err != nil {
		panic(err)
	}

	provider := sdkMetric.NewMeterProvider(
		sdkMetric.WithReader(
			sdkMetric.NewPeriodicReader(exporter)))
	return &otelMetrics{
		meter: provider.Meter("shorty"),
	}
}

type otelMetrics struct {
	meter metric.Meter
}

func (m otelMetrics) NewCounter(name, description string) Counter {
	result, err := m.meter.Int64Counter(name, metric.WithDescription(description))
	if err != nil {
		panic(err)
	}
	return otelCounter{result}
}

func (m otelMetrics) NewGauge(name, description string) Gauge {
	result, err := m.meter.Float64Gauge(name, metric.WithDescription(description))
	if err != nil {
		panic(err)
	}
	return otelGauge{result}
}

func (m otelMetrics) NewHistogram(name, description string, buckets []float64) Histogram {
	result, err := m.meter.Float64Histogram(
		name, metric.WithDescription(description),
		metric.WithExplicitBucketBoundaries(buckets...),
	)
	if err != nil {
		panic(err)
	}
	return otelHistogram{result}
}

type otelCounter struct {
	metric.Int64Counter
}

func (o otelCounter) Inc() {
	o.Int64Counter.Add(context.Background(), 1)
}

type otelGauge struct {
	metric.Float64Gauge
}

func (o otelGauge) Set(value float64) {
	o.Float64Gauge.Record(context.Background(), value)
}

type otelHistogram struct {
	metric.Float64Histogram
}

func (o otelHistogram) Observe(value float64) {
	o.Float64Histogram.Record(context.Background(), value)
}
