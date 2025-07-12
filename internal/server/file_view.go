package server

import (
	"fmt"
	"shorty/internal/server/pages"
	"shorty/internal/services/files"

	"github.com/gin-gonic/gin"
)

func (s *server) FileView(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		s.pages.NotFound(c)
		return
	}

	meta, err := s.FileService.GetFileMetadata(c, id)
	if err == files.ErrNotFound {
		s.pages.NotFound(c)
		return
	}
	if err != nil {
		s.pages.InternalError(c)
		return
	}

	captcha, _ := s.GuardService.CreateCaptcha(c)

	downloadUrl := fmt.Sprintf("%s/file/download/%s", s.Url, meta.Id)
	viewUrl := fmt.Sprintf("%s/file/view/%s", s.Url, meta.Id)

	s.pages.FileView(c, pages.FileViewParams{
		FileName:        meta.Name,
		FileSizeMB:      float32(meta.Size) / (1024 * 1024),
		FileViewUrl:     viewUrl,
		FileDownloadUrl: downloadUrl,
		CaptchaId:       captcha.Id,
		CaptchaBase64:   captcha.ImageBase64,
	})
}
