package server

import (
	"shorty/src/services/image"

	"github.com/gin-gonic/gin"
)

func (s *server) ImageResolve(c *gin.Context) {
	shortId := c.Param("id")
	if shortId == "" {
		s.pages.NotFound(c)
		return
	}

	img, err := s.ImageService.GetFile(c, shortId)
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
