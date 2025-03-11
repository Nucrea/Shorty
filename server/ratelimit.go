package server

import (
	"shorty/server/site"
	"shorty/src/services/ratelimit"

	"github.com/gin-gonic/gin"
)

func NewRatelimitM(
	ratelimitService *ratelimit.Service,
	site *site.Site,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := ratelimitService.Check(c, c.ClientIP())
		if err == ratelimit.ErrTooManyRequests {
			site.TooManyRequests(c)
			return
		}
		if err == ratelimit.ErrTemporaryBanned {
			site.TemporarilyBanned(c)
			return
		}

		c.Next()
	}
}
