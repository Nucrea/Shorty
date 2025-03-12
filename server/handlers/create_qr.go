package handlers

import (
	"shorty/server/site"
	"shorty/src/common/logger"
	"shorty/src/services/links"
	"shorty/src/services/ratelimit"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreateQRDeps struct {
	Log              logger.Logger
	Site             *site.Site
	LinkService      *links.Service
	RatelimitService *ratelimit.Service
}

func CreateQR(d CreateQRDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Query("url")
		if url == "" {
			// p.IndexPage.WithError(c, "Bad url")
			d.Site.CreateLink(c)
			return
		}

		qrCode, err := d.LinkService.CreateQR(c, url)
		if err == links.ErrBadUrl {
			d.Site.CreateLink(c)
			return
		}
		if err != nil {
			log.Error().Err(err).Msg("error creating qr code")
			d.Site.InternalError(c)
			return
		}

		d.Site.QRResult(c, qrCode)
	}
}
