package server

import (
	"net/url"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func (s *server) UserLogout(c *gin.Context) templ.Component {
	sessionIdCookie, err := c.Cookie(sessionCookieKey)
	if err != nil || sessionIdCookie == "" {
		c.Redirect(302, "/")
		return nil
	}

	err = s.UserService.DeleteSession(c, sessionIdCookie)
	if err != nil {
		return s.site.InternalError(c)
	}

	u, _ := url.Parse(s.Url)
	hostname := u.Hostname()

	c.SetCookie(sessionCookieKey, "", -1, "/", hostname, false, true)
	c.Redirect(302, "/")
	return nil
}
