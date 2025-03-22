package logging

import (
	"time"

	"github.com/rs/zerolog"
)

type Event interface {
	Str(key, val string) Event
	Int(key string, val int) Event
	Dur(key string, dur time.Duration) Event

	Err(err error) Event

	Send()
	Msg(msg string)
	Msgf(format string, v ...interface{})
}

type event struct {
	zerologEvent *zerolog.Event
}

func (e event) Int(key string, val int) Event {
	e.zerologEvent = e.zerologEvent.Int(key, val)
	return e
}

func (e event) Str(key, val string) Event {
	e.zerologEvent = e.zerologEvent.Str(key, val)
	return e
}

func (e event) Dur(key string, dur time.Duration) Event {
	e.zerologEvent = e.zerologEvent.Dur(key, dur)
	return e
}

func (e event) Err(err error) Event {
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
}

// Never call it twice: event gets disposed after first call
func (e event) Msgf(format string, v ...interface{}) {
	e.zerologEvent.Msgf(format, v...)
}
