package server

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *server) FileUpload(c *gin.Context) {
	header, err := c.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("error getting file from request")
		s.pages.InternalError(c)
		return
	}

	file, err := header.Open()
	if err != nil {
		log.Error().Err(err).Msg("error opening tmp file")
		s.pages.InternalError(c)
		return
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)

	info, err := s.FileService.UploadFile(c, header.Filename, bytes)
	if err != nil {
		log.Error().Err(err).Msg("error uploading file")
		s.pages.InternalError(c)
		return
	}

	viewUrl := fmt.Sprintf("/file/view/%s", info.ShortId)
	c.Redirect(302, viewUrl)
}
