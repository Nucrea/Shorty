package handlers

import (
	genericerror "shorty/server/pages/generic_error"
	"shorty/server/pages/index"
	"shorty/server/pages/result"
	"shorty/src/common/logger"
	"shorty/src/services/links"
	"shorty/src/services/ratelimit"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

type CreateHDeps struct {
	Log              logger.Logger
	IndexPage        *index.Page
	ResultPage       *result.Page
	ErrorPage        *genericerror.Page
	LinkService      *links.Service
	RatelimitService *ratelimit.Service
	Tracer           trace.Tracer
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
	if err == links.ErrBadUrl {
		p.IndexPage.WithError(c, "Bad url")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("error creating link")
		p.ErrorPage.InternalError(c)
		return
	}

	p.ResultPage.WithLink(c, link)
}

func createQrH(c *gin.Context, p CreateHDeps, url string) {
	qrCode, err := p.LinkService.CreateQR(c, url)
	if err == links.ErrBadUrl {
		p.IndexPage.WithError(c, "Bad url")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("error creating qr code")
		p.ErrorPage.InternalError(c)
		return
	}

	p.ResultPage.WithQRCode(c, qrCode)
}
