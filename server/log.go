package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func RequestLogM(log *zerolog.Logger) gin.HandlerFunc {
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

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		method := c.Request.Method
		statusCode := c.Writer.Status()

		log.Info().
			Str("request_id", requestId).
			Str("method", method).
			Str("path", path).
			Str("ip", c.ClientIP()).
			Int("status", statusCode).
			Dur("duration", duration).
			Send()
	}
}
