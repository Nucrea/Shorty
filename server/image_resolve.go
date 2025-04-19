package server

import (
	"fmt"
	"net/http"
	"shorty/src/common"
	"shorty/src/services/image"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func (s *server) ImageResolve(c *gin.Context) templ.Component {
	id := c.Param("id")
	if id == "" {
		return s.site.NotFound(c)
	}

	isThumbnail := false
	switch c.Param("type") {
	case "o":
		isThumbnail = false
	case "t":
		isThumbnail = true
	default:
		return s.site.NotFound(c)
	}

	meta, err := s.ImageService.GetImageMetadata(c, id)
	if err == image.ErrImageNotFound {
		return s.site.NotFound(c)
	}
	if err != nil {
		return s.site.InternalError(c)
	}

	if !isThumbnail {
		token, expiresStr := c.Query("token"), c.Query("expires")
		expires, _ := strconv.Atoi(expiresStr)

		expired := int(time.Now().Unix()) > expires
		valid := CheckResourceToken(id, int64(expires), token)

		if expired || !valid {
			s.Logger.WithContext(c).Info().Msgf("image (id=%s) token(%s) expired, redirecting to view", id, common.MaskSecret(token))
			viewUrl := fmt.Sprintf("%s/image/view/%s", s.Url, meta.Id)
			c.Redirect(302, viewUrl)
			return nil
		}
	}

	if oldEtag := c.GetHeader("If-None-Match"); oldEtag == meta.Hash {
		s.Logger.WithContext(c).Info().Msgf("image (id=%s) hash does not changed", id)
		c.AbortWithStatus(http.StatusNotModified)
		return nil
	}

	imageBytes, err := s.ImageService.GetImageBytes(c, id, isThumbnail)
	if err == image.ErrImageNotFound {
		return s.site.NotFound(c)
	}
	if err != nil {
		return s.site.InternalError(c)
	}

	c.Header("ETag", meta.Hash)
	c.Header("Cache-Control", "public, max-age=300")

	c.Data(200, "image/jpeg", imageBytes)
	return nil
}
