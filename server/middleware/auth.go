package middleware

import (
	"net/http"
	"shorty/src/services/users"

	"github.com/gin-gonic/gin"
)

const SessionCookieKey = "sessionId"
const SessionKey = "sessionObj"

func GetUserSession(c *gin.Context) *users.SessionDTO {
	session, ok := c.Get(SessionKey)
	if !ok {
		return nil
	}
	return session.(*users.SessionDTO)
}

func Authorization(userService *users.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionCookie, err := c.Request.Cookie(SessionCookieKey)
		if err != nil || sessionCookie == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		session, err := userService.Authorize(c, sessionCookie.Value)
		if err != nil || session == nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set(SessionKey, session)
		c.Next()
	}
}
