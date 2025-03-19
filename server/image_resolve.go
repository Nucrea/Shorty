package server

import (
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

	img, err := s.ImageService.GetImageBytes(c, id, isThumbnail)
	if err == image.ErrImageNotFound {
		s.pages.NotFound(c)
		return
	}
	if err != nil {
		s.pages.InternalError(c)
		return
	}

	c.Data(200, "image/jpeg", img)
}
