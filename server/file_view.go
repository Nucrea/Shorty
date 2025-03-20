package server

import (
	"fmt"
	"shorty/server/pages"
	"shorty/src/services/files"
	"time"

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

	expiresAt := time.Now().Add(10 * time.Second).UnixMicro()
	token := s.GuardService.CreateResourceToken(meta.Id, expiresAt)

	downloadUrl := fmt.Sprintf("%s/f/%s/%s?token=%s&expires=%d", s.Url, meta.Id, meta.Name, token, expiresAt)
	viewUrl := fmt.Sprintf("%s/file/view/%s", s.Url, meta.Id)

	s.pages.FileView(c, pages.ViewFileParams{
		FileName:        meta.Name,
		FileSizeMB:      float32(meta.Size) / (1024 * 1024),
		FileViewUrl:     viewUrl,
		FileDownloadUrl: downloadUrl,
		CaptchaId:       captcha.Id,
		CaptchaBase64:   captcha.ImageBase64,
	})
}
