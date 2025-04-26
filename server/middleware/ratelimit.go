package middleware

import (
	"net/http"
	"shorty/src/services/guard"

	"github.com/gin-gonic/gin"
)

func Ratelimit(guardService *guard.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := guardService.CheckIP(c, c.ClientIP())
		if err == guard.ErrTooManyRequests {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		if err == guard.ErrTemporaryBanned {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}
