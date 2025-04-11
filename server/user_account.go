package server

import (
	"github.com/gin-gonic/gin"
)

func (s *server) UserAccount(c *gin.Context) {
	sessionId, err := c.Cookie("sessionId")
	if err != nil {
		c.Redirect(302, "login")
		return
	}

	_, err = s.UserService.Authorize(c, sessionId)
	if err != nil {
		c.Redirect(302, "login")
		return
	}

	s.pages.AccountView(c)
}
