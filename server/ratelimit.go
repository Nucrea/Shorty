package server

import (
	genericerror "shorty/server/pages/generic_error"
	"shorty/src/services/ratelimit"

	"github.com/gin-gonic/gin"
)

func NewRatelimitM(
	ratelimitService *ratelimit.Service,
	errorPage *genericerror.Page,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := ratelimitService.Check(c, c.ClientIP())
		if err == ratelimit.ErrTooManyRequests {
			errorPage.TooMuchRequests(c)
			return
		}
		if err == ratelimit.ErrTemporaryBanned {
			errorPage.TemporarilyBanned(c)
			return
		}

		c.Next()
	}
}
