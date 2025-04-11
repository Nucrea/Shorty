package server

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (s *server) UserLogin(c *gin.Context) {
	email, password := c.PostForm("email"), c.PostForm("password")

	session, err := s.UserService.Login(c, email, password)
	if err != nil {
		url := fmt.Sprintf("/login?err=%s", err.Error())
		c.Redirect(302, url)
		return
	}

	u, _ := url.Parse(s.Url)
	hostname := u.Hostname()

	c.SetCookie(sessionCookieKey, session.Token, 3600, "/", hostname, false, true)
	c.Redirect(302, "/account")
}
