package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const RequestIdKey = "logger_request_id"

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

type logger struct {
	service    string
	requestId  string
	zeroLogger *zerolog.Logger
}

func (l logger) Log() Event {
	return l.wrapEvent(l.zeroLogger.Log())
}

func (l logger) Info() Event {
	return l.wrapEvent(l.zeroLogger.Info())
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
	}
}

func (l logger) WithService(service string) Logger {
	return logger{
		service:    service,
		requestId:  l.requestId,
		zeroLogger: l.zeroLogger,
	}
}

func SetCtxRequestId(ginCtx *gin.Context, requestId string) {
	ginCtx.Set(RequestIdKey, requestId)
}
