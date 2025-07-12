package logging

type Option interface {
	apply(config) config
}

type optionFunc func(config) config

func (fn optionFunc) apply(c config) config {
	return fn(c)
}

type config struct {
	OutputFile string
	OtelUrl    string
}

func WithFile(path string) Option {
	return optionFunc(func(c config) config {
		c.OutputFile = path
		return c
	})
}

func WithOpenTelemetry(otelUrl string) Option {
	return optionFunc(func(c config) config {
		c.OtelUrl = otelUrl
		return c
	})
}
