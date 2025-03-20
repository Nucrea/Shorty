package server

import (
	"fmt"
	"net/url"
	"shorty/src/services/files"

	"github.com/gin-gonic/gin"
)

func (s *server) FileResolve(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		s.pages.NotFound(c)
		return
	}

	captchaId, captchaToken := c.Query("id"), c.Query("token")
	err := s.GuardService.CheckCaptcha(c, captchaId, captchaToken)
	if err != nil {
		c.Redirect(302, fmt.Sprintf("/file/view/%s?err=%s", id, url.QueryEscape("captcha wrong or expired")))
		return
	}

	fileBytes, err := s.FileService.GetFileBytes(c, id)
	if err == files.ErrNotFound {
		s.pages.NotFound(c)
		return
	}
	if err != nil {
		s.pages.InternalError(c)
		return
	}

	c.Data(200, "application/octet-stream", fileBytes)
}
