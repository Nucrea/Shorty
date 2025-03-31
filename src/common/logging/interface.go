package logging

import (
	"context"
	"time"
)

type Logger interface {
	Log() Event
	Info() Event
	Warning() Event
	Error() Event
	Fatal() Event

	Printf(format string, v ...any)

	WithContext(ctx context.Context) Logger
	WithService(name string) Logger
}

type Event interface {
	Str(key, val string) Event
	Int(key string, val int) Event
	Dur(key string, dur time.Duration) Event

	Err(err error) Event

	Send()
	Msg(msg string)
	Msgf(format string, v ...interface{})
}
