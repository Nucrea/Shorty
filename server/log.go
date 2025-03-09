package server

import (
	"shorty/src/common/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestLogM(log logger.Logger) gin.HandlerFunc {
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

		log.Info().
			Str("requestId", requestId).
			Str("method", method).
			Str("path", path).
			Str("ip", c.ClientIP()).
			Int("status", statusCode).
			Dur("duration", duration).
			Send()
	}
}
