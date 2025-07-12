package logging

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const RequestIdKey = "logger_request_id"

func SetCtxRequestId(ginCtx *gin.Context, requestId string) {
	ginCtx.Set(RequestIdKey, requestId)
}

func Fatal(err error) {
	l := zerolog.New(os.Stdout).
		Level(zerolog.DebugLevel).
		With().Timestamp().
		Logger()
	l.Fatal().Err(err).Msg("global fatal error occured")
}
