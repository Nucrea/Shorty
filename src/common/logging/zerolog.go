package logging

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewZerolog(outputFile string) (Logger, error) {
	writers := []io.Writer{}
	writers = append(writers, os.Stdout)

	if outputFile != "" {
		file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}

	writer := io.MultiWriter(writers...)
	l := zerolog.New(writer).
		Level(zerolog.DebugLevel).
		With().Timestamp().
		Logger()

	return &zerologger{
		zeroLogger: &l,
	}, nil
}

type zerologger struct {
	service    string
	requestId  string
	zeroLogger *zerolog.Logger
}

func (l zerologger) Log() Event {
	return l.wrapEvent(l.zeroLogger.Log())
}

func (l zerologger) Info() Event {
	return l.wrapEvent(l.zeroLogger.Info())
}

func (l zerologger) Warning() Event {
	return l.wrapEvent(l.zeroLogger.Warn())
}

func (l zerologger) Error() Event {
	return l.wrapEvent(l.zeroLogger.Error())
}

func (l zerologger) Fatal() Event {
	return l.wrapEvent(l.zeroLogger.Fatal())
}

func (l zerologger) Printf(format string, v ...any) {
	l.zeroLogger.Printf(format, v...)
}

func (l zerologger) wrapEvent(zerologEvent *zerolog.Event) Event {
	var e Event = zevent{zerologEvent}

	if l.requestId != "" {
		e = e.Str("requestId", l.requestId)
	}
	if l.service != "" {
		e = e.Str("service", l.service)
	}

	return e
}

func (l zerologger) WithContext(ctx context.Context) Logger {
	requestIdVal := ctx.Value(RequestIdKey)
	requestId, ok := requestIdVal.(string)
	if !ok || requestId == "" {
		return l
	}

	return zerologger{
		service:    l.service,
		requestId:  requestId,
		zeroLogger: l.zeroLogger,
	}
}

func (l zerologger) WithService(service string) Logger {
	return zerologger{
		service:    service,
		requestId:  l.requestId,
		zeroLogger: l.zeroLogger,
	}
}

type zevent struct {
	zerologEvent *zerolog.Event
}

func (e zevent) Int(key string, val int) Event {
	e.zerologEvent = e.zerologEvent.Int(key, val)
	return e
}

func (e zevent) Str(key, val string) Event {
	e.zerologEvent = e.zerologEvent.Str(key, val)
	return e
}

func (e zevent) Dur(key string, dur time.Duration) Event {
	e.zerologEvent = e.zerologEvent.Dur(key, dur)
	return e
}

func (e zevent) Err(err error) Event {
	e.zerologEvent = e.zerologEvent.Err(err)
	return e
}

// Never call it twice: event gets disposed after first call
func (e zevent) Send() {
	e.zerologEvent.Send()
}

// Never call it twice: event gets disposed after first call
func (e zevent) Msg(msg string) {
	e.zerologEvent.Msg(msg)
}

// Never call it twice: event gets disposed after first call
func (e zevent) Msgf(format string, v ...interface{}) {
	e.zerologEvent.Msgf(format, v...)
}
