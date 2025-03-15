package server

import (
	"fmt"
	"shorty/server/pages"
	"shorty/src/services/image"

	"github.com/gin-gonic/gin"
)

func (s *server) ImageView(c *gin.Context) {
	shortId := c.Param("id")
	if shortId == "" {
		s.pages.NotFound(c)
		return
	}

	info, err := s.ImageService.GetImageInfo(c, shortId)
	if err == image.ErrImageNotFound {
		s.pages.NotFound(c)
		return
	}
	if err != nil {
		s.pages.InternalError(c)
		return
	}

	viewUrl := fmt.Sprintf("%s/image/view/%s", s.Url, info.ShortId)
	thumbUrl := fmt.Sprintf("%s/i/t/%s", s.Url, info.ThumbnailId)
	imgUrl := fmt.Sprintf("%s/i/f/%s", s.Url, info.ImageId)

	s.pages.ImageView(c, pages.ViewImageParams{
		FileName:     info.Name,
		SizeMB:       float32(info.Size) / (1024 * 1024),
		ViewUrl:      viewUrl,
		ImageUrl:     imgUrl,
		ThumbnailUrl: thumbUrl,
	})
}
