package server

import (
	"fmt"
	"net/url"
	"shorty/src/services/files"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *server) FileDownload(c *gin.Context) {
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

	meta, err := s.FileService.GetFileMetadata(c, id)
	if err == files.ErrNotFound {
		s.pages.NotFound(c)
		return
	}
	if err != nil {
		s.pages.InternalError(c)
		return
	}

	token := NewResourceToken(id, time.Now().Add(15*time.Minute))
	fileRawUrl := fmt.Sprintf("%s/f/%s/%s?token=%s&expires=%d", s.Url, meta.Id, meta.Name, token.Value, token.Exipres)

	s.pages.FileDownload(c, fileRawUrl)
}
