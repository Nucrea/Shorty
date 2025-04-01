package logging

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(opts ...Option) (Logger, error) {
	conf := config{}
	for _, opt := range opts {
		conf = opt.apply(conf)
	}

	var writer io.Writer = os.Stdout

	if conf.OutputFile != "" {
		file, err := os.OpenFile(conf.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			return nil, err
		}

		writer = io.MultiWriter(os.Stdout, file)
	}

	// otel := newOtelNoop()
	// if conf.OtelUrl != "" {
	// 	otel = newOtel(conf.OtelUrl)
	// }

	l := zerolog.New(writer).
		Level(zerolog.DebugLevel).
		With().Timestamp().
		Logger()

	return &logger{
		zeroLogger: &l,
		// otel:       otel,
	}, nil
}

type logger struct {
	service    string
	requestId  string
	zeroLogger *zerolog.Logger
	// otel       otellog.Logger
}

func (l logger) Log() Event {
	return l.wrapEvent(l.zeroLogger.Log())
}

func (l logger) Info() Event {
	return l.wrapEvent(l.zeroLogger.Info())
}

func (l logger) Debug() Event {
	return l.wrapEvent(l.zeroLogger.Debug())
}

func (l logger) Warning() Event {
	return l.wrapEvent(l.zeroLogger.Warn())
}

func (l logger) Error() Event {
	return l.wrapEvent(l.zeroLogger.Error())
}

func (l logger) Fatal() Event {
	return l.wrapEvent(l.zeroLogger.Fatal())
}

func (l logger) Printf(format string, v ...any) {
	l.zeroLogger.Printf(format, v...)
}

func (l logger) wrapEvent(zerologEvent *zerolog.Event) Event {
	var e Event = event{zerologEvent}

	if l.requestId != "" {
		e = e.Str("requestId", l.requestId)
	}
	if l.service != "" {
		e = e.Str("service", l.service)
	}

	return e
}

func (l logger) WithContext(ctx context.Context) Logger {
	requestIdVal := ctx.Value(RequestIdKey)
	requestId, ok := requestIdVal.(string)
	if !ok || requestId == "" {
		return l
	}

	return logger{
		service:    l.service,
		requestId:  requestId,
		zeroLogger: l.zeroLogger,
		// otel:       l.otel,
	}
}

func (l logger) WithService(service string) Logger {
	return logger{
		service:    service,
		requestId:  l.requestId,
		zeroLogger: l.zeroLogger,
		// otel:       l.otel,
	}
}

type event struct {
	zerologEvent *zerolog.Event
	// otel         otellog.Logger
	// record       otellog.Record
}

func (e event) Int(key string, val int) Event {
	// (&e.record).AddAttributes(otellog.Int(key, val))
	e.zerologEvent = e.zerologEvent.Int(key, val)
	return e
}

func (e event) Str(key, val string) Event {
	// (&e.record).AddAttributes(otellog.String(key, val))
	e.zerologEvent = e.zerologEvent.Str(key, val)
	return e
}

func (e event) Dur(key string, dur time.Duration) Event {
	// (&e.record).AddAttributes(otellog.String(key, dur.String()))
	e.zerologEvent = e.zerologEvent.Dur(key, dur)
	return e
}

func (e event) Err(err error) Event {
	// (&e.record).AddAttributes(otellog.String("error", err.Error()))
	e.zerologEvent = e.zerologEvent.Err(err)
	return e
}

// Never call it twice: event gets disposed after first call
func (e event) Send() {
	e.zerologEvent.Send()
}

// Never call it twice: event gets disposed after first call
func (e event) Msg(msg string) {
	e.zerologEvent.Msg(msg)

	// r := &e.record
	// r.AddAttributes(otellog.String("message", msg))
	// r.SetSeverity(otellog.SeverityInfo)
	// r.SetBody(otellog.StringValue(msg))
	// e.otel.Emit(context.Background(), e.record)
}

// Never call it twice: event gets disposed after first call
func (e event) Msgf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	e.zerologEvent.Msg(msg)

	// r := &e.record
	// r.AddAttributes(otellog.String("message", msg))
	// r.SetSeverity(otellog.SeverityInfo)
	// r.SetBody(otellog.StringValue(msg))
	// e.otel.Emit(context.Background(), e.record)
}
