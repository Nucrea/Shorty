package server

import (
	"fmt"
	"io"
	"net/url"
	"shorty/src/services/image"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *server) ImageUpload(c *gin.Context) templ.Component {
	id, token := c.PostForm("id"), c.PostForm("token")
	err := s.GuardService.CheckCaptcha(c, id, token)
	if err != nil {
		c.Redirect(302, "/image?err="+url.QueryEscape("captcha wrong or expired"))
		return nil
	}

	header, err := c.FormFile("image")
	if err != nil {
		log.Error().Err(err).Msg("error getting image from request")
		c.Redirect(302, "/image?err="+url.QueryEscape(err.Error()))
		return nil
	}

	file, err := header.Open()
	if err != nil {
		log.Error().Err(err).Msg("error opening image file")
		return s.site.InternalError(c)
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)

	meta, err := s.ImageService.UploadImage(c, header.Filename, bytes)
	if err == image.ErrInvalidFormat || err == image.ErrUnsupportedFormat || err == image.ErrImageTooLarge {
		log.Error().Err(err).Msg("error getting image from request")
		c.Redirect(302, "/image?err="+url.QueryEscape(err.Error()))
		return nil
	}
	if err != nil {
		log.Error().Err(err).Msg("error creating image")
		return s.site.InternalError(c)
	}

	imgUrl := fmt.Sprintf("/image/view/%s", meta.Id)
	c.Redirect(302, imgUrl)
	return nil
}
