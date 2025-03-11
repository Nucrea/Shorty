package handlers

import (
	"fmt"
	"io"
	"shorty/server/site"
	"shorty/src/common/logger"
	"shorty/src/services/image"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UploadImageDeps struct {
	BaseUrl      string
	Log          logger.Logger
	Site         *site.Site
	ImageService *image.Service
}

func UploadImage(d UploadImageDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		header, err := c.FormFile("image")
		if err != nil {
			log.Error().Err(err).Msg("error getting image from request")
			d.Site.InternalError(c)
			return
		}

		file, err := header.Open()
		if err != nil {
			log.Error().Err(err).Msg("error opening image file")
			d.Site.InternalError(c)
			return
		}
		defer file.Close()

		bytes, _ := io.ReadAll(file)

		info, err := d.ImageService.CreateImage(c, header.Filename, bytes)
		if err != nil {
			log.Error().Err(err).Msg("error creating image")
			d.Site.InternalError(c)
			return
		}

		imgUrl := fmt.Sprintf("%s/image/view/%s", d.BaseUrl, info.ShortId)
		c.Redirect(302, imgUrl)
	}
}
