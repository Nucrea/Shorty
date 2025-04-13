package middleware

import (
	"shorty/server/site"
	"shorty/src/services/guard"

	"github.com/gin-gonic/gin"
)

func Ratelimit(
	guardService *guard.Service,
	site *site.Site,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := guardService.CheckIP(c, c.ClientIP())
		if err == guard.ErrTooManyRequests {
			site.TooManyRequests(c)
			return
		}
		if err == guard.ErrTemporaryBanned {
			site.TemporarilyBanned(c)
			return
		}

		c.Next()
	}
}
