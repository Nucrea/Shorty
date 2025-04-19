package server

import (
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func (s *server) ImageForm(c *gin.Context) templ.Component {
	captcha, _ := s.GuardService.CreateCaptcha(c)
	return s.site.ImageForm(c, captcha.Id, captcha.ImageBase64)
}
