package server

import (
	"fmt"
	"net/url"
	"shorty/src/common"
	"shorty/src/services/files"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func (s *server) FileResolve(c *gin.Context) templ.Component {
	id := c.Param("id")
	if id == "" {
		return s.site.NotFound(c)
	}

	token, expiresStr := c.Query("token"), c.Query("expires")
	expires, _ := strconv.Atoi(expiresStr)

	expired := int(time.Now().Unix()) > expires
	valid := CheckResourceToken(id, int64(expires), token)

	if expired || !valid {
		s.Logger.WithContext(c).Info().Msgf("file (id=%s) token(%s) expired, redirecting to view", id, common.MaskSecret(token))
		viewUrl := fmt.Sprintf("%s/file/view/%s?err=%s", s.Url, id, url.QueryEscape("Download link expired"))
		c.Redirect(302, viewUrl)
		return nil
	}

	fileBytes, err := s.FileService.GetFileBytes(c, id)
	if err == files.ErrNotFound {
		return s.site.NotFound(c)
	}
	if err != nil {
		return s.site.InternalError(c)
	}

	c.Data(200, "application/octet-stream", fileBytes)
	return nil
}
