package handlers

import (
	"io"
	genericerror "shorty/server/pages/generic_error"
	"shorty/server/pages/index"
	"shorty/server/pages/result"
	"shorty/server/pages/upload"
	"shorty/src/common/logger"
	"shorty/src/services/image"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ImageHDeps struct {
	Log          logger.Logger
	ImageService *image.Service
	IndexPage    *index.Page
	ResultPage   *result.Page
	ErrorPage    *genericerror.Page
	UploadPage   *upload.Page
}

func NewImageCreateH(p ImageHDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		header, err := c.FormFile("image")
		if err != nil {
			log.Error().Err(err).Msg("error getting image from request")
			p.ErrorPage.InternalError(c)
			return
		}

		file, err := header.Open()
		if err != nil {
			log.Error().Err(err).Msg("error opening image file")
			p.ErrorPage.InternalError(c)
			return
		}
		defer file.Close()

		bytes, _ := io.ReadAll(file)

		info, err := p.ImageService.CreateImage(c, header.Filename, bytes)
		if err != nil {
			log.Error().Err(err).Msg("error creating image")
			p.ErrorPage.InternalError(c)
			return
		}

		p.ResultPage.WithLink(c, info.ShortId)
	}
}
