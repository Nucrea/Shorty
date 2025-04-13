package server

import (
	"fmt"
	"shorty/server/site/pages"
	"shorty/src/services/image"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *server) ImageView(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		s.site.NotFound(c)
		return
	}

	year, month, day := time.Now().Add(24 * time.Hour).Date()
	expiresAt := time.Date(year, month, day, 5, 0, 0, 0, time.Now().Location())

	meta, err := s.ImageService.GetImageMetadata(c, id)
	if err == image.ErrImageNotFound {
		s.site.NotFound(c)
		return
	}
	if err != nil {
		s.site.InternalError(c)
		return
	}

	viewUrl := fmt.Sprintf("%s/image/view/%s", s.Url, meta.Id)
	thumbUrl := fmt.Sprintf("%s/i/t/%s", s.Url, meta.Id)

	token := NewResourceToken(meta.Id, expiresAt)
	imgUrl := fmt.Sprintf("%s/i/o/%s?token=%s&expires=%d", s.Url, meta.Id, token.Value, token.Exipres)

	s.site.ImageView(c, pages.ImageViewParams{
		FileName:     meta.Name,
		SizeMB:       float32(meta.Size) / (1024 * 1024),
		ViewUrl:      viewUrl,
		ImageUrl:     imgUrl,
		ThumbnailUrl: thumbUrl,
	})
}
