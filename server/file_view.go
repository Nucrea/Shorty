package server

import (
	"fmt"
	"shorty/server/site/pages"
	"shorty/src/services/files"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func (s *server) FileView(c *gin.Context) templ.Component {
	id := c.Param("id")
	if id == "" {
		return s.site.NotFound(c)
	}

	meta, err := s.FileService.GetFileMetadata(c, id)
	if err == files.ErrNotFound {
		return s.site.NotFound(c)
	}
	if err != nil {
		return s.site.InternalError(c)
	}

	captcha, _ := s.GuardService.CreateCaptcha(c)

	downloadUrl := fmt.Sprintf("%s/file/download/%s", s.Url, meta.Id)
	viewUrl := fmt.Sprintf("%s/file/view/%s", s.Url, meta.Id)

	return s.site.FileView(c, pages.FileViewParams{
		FileName:        meta.Name,
		FileSizeMB:      float32(meta.Size) / (1024 * 1024),
		FileViewUrl:     viewUrl,
		FileDownloadUrl: downloadUrl,
		CaptchaId:       captcha.Id,
		CaptchaBase64:   captcha.ImageBase64,
	})
}
