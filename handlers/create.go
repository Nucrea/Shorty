package handlers

import (
	genericerror "shorty/pages/generic_error"
	"shorty/pages/index"
	"shorty/pages/result"
	"shorty/services/links"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CreateHDeps struct {
	Log         *zerolog.Logger
	IndexPage   *index.Page
	ResultPage  *result.Page
	ErrorPage   *genericerror.Page
	LinkService *links.Service
}

func NewLinkCreateH(p CreateHDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Query("url")
		if url == "" {
			p.IndexPage.WithError(c, "Bad url")
			return
		}

		if c.Query("qr") == "on" {
			createQrH(c, p, url)
		} else {
			createLinkH(c, p, url)
		}
	}
}

func createLinkH(c *gin.Context, p CreateHDeps, url string) {
	link, err := p.LinkService.CreateLink(c, url)
	if err != nil {
		log.Error().Err(err).Msg("error creating link")
		p.ErrorPage.InternalError(c)
		return
	}

	p.ResultPage.WithLink(c, link)
}

func createQrH(c *gin.Context, p CreateHDeps, url string) {
	qrCode, err := p.LinkService.CreateQR(c, url)
	if err != nil {
		log.Error().Err(err).Msg("error creating qr code")
		p.ErrorPage.InternalError(c)
		return
	}

	p.ResultPage.WithQRCode(c, qrCode)
}
