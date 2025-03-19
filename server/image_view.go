package server

import (
	"fmt"
	"shorty/server/pages"
	"shorty/src/services/image"

	"github.com/gin-gonic/gin"
)

func (s *server) ImageView(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		s.pages.NotFound(c)
		return
	}

	meta, err := s.ImageService.GetImageMetadata(c, id)
	if err == image.ErrImageNotFound {
		s.pages.NotFound(c)
		return
	}
	if err != nil {
		s.pages.InternalError(c)
		return
	}

	viewUrl := fmt.Sprintf("%s/image/view/%s", s.Url, meta.Id)
	thumbUrl := fmt.Sprintf("%s/i/t/%s", s.Url, meta.Id)
	imgUrl := fmt.Sprintf("%s/i/o/%s", s.Url, meta.Id)

	s.pages.ImageView(c, pages.ViewImageParams{
		FileName:     meta.Name,
		SizeMB:       float32(meta.Size) / (1024 * 1024),
		ViewUrl:      viewUrl,
		ImageUrl:     imgUrl,
		ThumbnailUrl: thumbUrl,
	})
}
