package server

import (
	"fmt"
	"shorty/server/pages"
	"shorty/src/services/image"
	"time"

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

	expiresAt := time.Now().Add(time.Hour).UnixMicro()
	token := s.GuardService.CreateResourceToken(meta.Id, expiresAt)
	imgUrl := fmt.Sprintf("%s/i/o/%s?token=%s&expires=%d", s.Url, meta.Id, token, expiresAt)

	s.pages.ImageView(c, pages.ViewImageParams{
		FileName:     meta.Name,
		SizeMB:       float32(meta.Size) / (1024 * 1024),
		ViewUrl:      viewUrl,
		ImageUrl:     imgUrl,
		ThumbnailUrl: thumbUrl,
	})
}
