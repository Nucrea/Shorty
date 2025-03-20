package server

import (
	"fmt"
	"net/url"
	"shorty/src/common"
	"shorty/src/services/files"
	"strconv"
	"time"

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

	token, expiresStr := c.Query("token"), c.Query("expires")
	expires, _ := strconv.Atoi(expiresStr)

	expired := int(time.Now().UnixMicro()) > expires
	valid := s.GuardService.CheckResourceToken(id, int64(expires), token)

	if expired || !valid {
		s.Log.WithContext(c).Info().Msgf("file (id=%s) token(%s) expired, redirecting to view", id, common.MaskSecret(token))
		viewUrl := fmt.Sprintf("%s/file/view/%s", s.Url, id)
		c.Redirect(302, viewUrl)
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
