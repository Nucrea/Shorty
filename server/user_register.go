package server

import (
	"fmt"
	"shorty/src/services/users"

	"github.com/gin-gonic/gin"
)

func (s *server) UserRegister(c *gin.Context) {
	email, password := c.PostForm("email"), c.PostForm("password")

	_, err := s.UserService.Create(c, users.CreateUserParams{
		Email:    email,
		Password: password,
	})
	if err != nil {
		url := fmt.Sprintf("/register?err=%s", err.Error())
		c.Redirect(302, url)
		return
	}

	session, err := s.UserService.Login(c, email, password)
	if err != nil {
		url := fmt.Sprintf("/register?err=%s", err.Error())
		c.Redirect(302, url)
		return
	}

	c.SetCookie("sessionId", session.Token, 3600, "/", s.Url, false, true)
	c.Redirect(302, "/account")
}
