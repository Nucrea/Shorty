package server

import (
	"github.com/gin-gonic/gin"
)

func (s *server) ImageForm(c *gin.Context) {
	captcha, _ := s.GuardService.CreateCaptcha(c)
	s.site.ImageForm(c, captcha.Id, captcha.ImageBase64)
}
