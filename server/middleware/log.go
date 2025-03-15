package middleware

import (
	"shorty/src/common/logger"
	"shorty/src/common/tracing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Log(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetHeader("X-Request-Id")
		if requestId == "" {
			requestId = uuid.New().String()
		}
		c.Header("X-Request-Id", requestId)

		path := c.Request.URL.Path
		// if c.Request.URL.RawQuery != "" {
		// 	path = path + "?" + c.Request.URL.RawQuery
		// }

		logger.SetCtxRequestId(c, requestId)

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		method := c.Request.Method
		statusCode := c.Writer.Status()

		traceId := c.GetString(tracing.TraceIdCtxKey)

		info := log.Info().Str("requestId", requestId)
		if traceId != "" {
			info = info.Str("traceId", traceId)
		}

		info.Str("method", method).
			Str("path", path).
			Str("ip", c.ClientIP()).
			Int("status", statusCode).
			Dur("duration", duration).
			Send()
	}
}
