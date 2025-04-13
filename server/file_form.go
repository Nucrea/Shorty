package server

import (
	"github.com/gin-gonic/gin"
)

func (s *server) FileForm(c *gin.Context) {
	captcha, _ := s.GuardService.CreateCaptcha(c)
	s.site.FileForm(c, captcha.Id, captcha.ImageBase64)
}
