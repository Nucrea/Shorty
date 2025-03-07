package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RequestLogM(log *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// requestId := c.GetHeader("X-Request-Id")
		// if requestId == "" {
		// 	requestId = uuid.New().String()
		// }
		// c.Header("X-Request-Id", requestId)

		path := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			path = path + "?" + c.Request.URL.RawQuery
		}

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		method := c.Request.Method
		statusCode := c.Writer.Status()

		log.Info().Msgf("%s %s %d %v", method, path, statusCode, latency)
	}
}
