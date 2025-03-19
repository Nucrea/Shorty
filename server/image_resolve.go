package server

import (
	"net/http"
	"shorty/src/services/image"

	"github.com/gin-gonic/gin"
)

func (s *server) ImageResolve(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		s.pages.NotFound(c)
		return
	}

	isThumbnail := false
	switch c.Param("type") {
	case "o":
		isThumbnail = false
	case "t":
		isThumbnail = true
	default:
		s.pages.NotFound(c)
		return
	}

	meta, err := s.ImageService.GetImageMetadata(c, id)
	if err != nil {
		s.pages.InternalError(c)
		return
	}

	if oldEtag := c.GetHeader("If-None-Match"); oldEtag == meta.Hash {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}

	imageBytes, err := s.ImageService.GetImageBytes(c, id, isThumbnail)
	if err == image.ErrImageNotFound {
		s.pages.NotFound(c)
		return
	}
	if err != nil {
		s.pages.InternalError(c)
		return
	}

	c.Header("ETag", meta.Hash)
	c.Header("Cache-Control", "public, max-age=300")

	c.Data(200, "image/jpeg", imageBytes)
}
