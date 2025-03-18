package server

import (
	"fmt"
	"io"
	"net/url"
	"shorty/src/services/image"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *server) ImageUpload(c *gin.Context) {
	header, err := c.FormFile("image")
	if err != nil {
		log.Error().Err(err).Msg("error getting image from request")
		c.Redirect(302, "image?err="+url.QueryEscape(err.Error()))
		return
	}

	file, err := header.Open()
	if err != nil {
		log.Error().Err(err).Msg("error opening image file")
		s.pages.InternalError(c)
		return
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)

	info, err := s.ImageService.UploadImage(c, header.Filename, bytes)
	if err == image.ErrInvalidFormat || err == image.ErrUnsupportedFormat || err == image.ErrImageTooLarge {
		log.Error().Err(err).Msg("error getting image from request")
		c.Redirect(302, "image?err="+url.QueryEscape(err.Error()))
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("error creating image")
		s.pages.InternalError(c)
		return
	}

	imgUrl := fmt.Sprintf("/image/view/%s", info.ShortId)
	c.Redirect(302, imgUrl)
}
