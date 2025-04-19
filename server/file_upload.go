package server

import (
	"fmt"
	"io"
	"net/url"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *server) FileUpload(c *gin.Context) templ.Component {
	id, token := c.PostForm("id"), c.PostForm("token")
	err := s.GuardService.CheckCaptcha(c, id, token)
	if err != nil {
		c.Redirect(302, "/file?err="+url.QueryEscape("captcha wrong or expired"))
		return nil
	}

	header, err := c.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("error getting file from request")
		return s.site.InternalError(c)
	}

	file, err := header.Open()
	if err != nil {
		log.Error().Err(err).Msg("error opening tmp file")
		return s.site.InternalError(c)
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)

	meta, err := s.FileService.UploadFile(c, header.Filename, bytes)
	if err != nil {
		log.Error().Err(err).Msg("error uploading file")
		return s.site.InternalError(c)
	}

	viewUrl := fmt.Sprintf("/file/view/%s", meta.Id)
	c.Redirect(302, viewUrl)
	return nil
}
