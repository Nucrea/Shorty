package server

import (
	"fmt"
	"net/url"
	"shorty/internal/common"
	"shorty/internal/services/files"
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

	token, expiresStr := c.Query("token"), c.Query("expires")
	expires, _ := strconv.Atoi(expiresStr)

	expired := int(time.Now().Unix()) > expires
	valid := CheckResourceToken(id, int64(expires), token)

	if expired || !valid {
		s.Logger.WithContext(c).Info().Msgf("file (id=%s) token(%s) expired, redirecting to view", id, common.MaskSecret(token))
		viewUrl := fmt.Sprintf("%s/file/view/%s?err=%s", s.Url, id, url.QueryEscape("Download link expired"))
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
