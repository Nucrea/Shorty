package server

import (
	"fmt"
	"shorty/server/pages"
	"shorty/src/services/files"

	"github.com/gin-gonic/gin"
)

func (s *server) FileView(c *gin.Context) {
	shortId := c.Param("id")
	if shortId == "" {
		s.pages.NotFound(c)
		return
	}

	info, err := s.FileService.GetFileInfo(c, shortId)
	if err == files.ErrNotFound {
		s.pages.NotFound(c)
		return
	}
	if err != nil {
		s.pages.InternalError(c)
		return
	}

	viewUrl := fmt.Sprintf("%s/file/view/%s", s.Url, info.ShortId)
	downloadUrl := fmt.Sprintf("%s/f/%s/%s", s.Url, info.ResourceId, info.Name)

	s.pages.FileView(c, pages.ViewFileParams{
		FileName:        info.Name,
		FileSizeMB:      float32(info.Size) / (1024 * 1024),
		FileViewUrl:     viewUrl,
		FileDownloadUrl: downloadUrl,
	})
}
