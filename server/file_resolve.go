package server

import (
	"shorty/src/services/files"

	"github.com/gin-gonic/gin"
)

func (s *server) FileResolve(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		s.pages.NotFound(c)
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
