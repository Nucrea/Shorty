package server

import (
	"fmt"
	"net/url"
	"shorty/src/services/files"
	"time"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func (s *server) FileDownload(c *gin.Context) templ.Component {
	id := c.Param("id")
	if id == "" {
		return s.site.NotFound(c)
	}

	captchaId, captchaToken := c.Query("id"), c.Query("token")
	err := s.GuardService.CheckCaptcha(c, captchaId, captchaToken)
	if err != nil {
		c.Redirect(302, fmt.Sprintf("/file/view/%s?err=%s", id, url.QueryEscape("captcha wrong or expired")))
		return nil
	}

	meta, err := s.FileService.GetFileMetadata(c, id)
	if err == files.ErrNotFound {
		return s.site.NotFound(c)
	}
	if err != nil {
		return s.site.InternalError(c)
	}

	token := NewResourceToken(id, time.Now().Add(15*time.Minute))
	fileRawUrl := fmt.Sprintf("%s/f/%s/%s?token=%s&expires=%d", s.Url, meta.Id, meta.Name, token.Value, token.Exipres)

	return s.site.FileDownload(c, fileRawUrl)
}
