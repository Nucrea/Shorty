package handlers

import (
	"shorty/server/site"
	"shorty/src/common/logger"
	"shorty/src/services/links"
	"shorty/src/services/ratelimit"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreateLinkDeps struct {
	Log              logger.Logger
	Site             *site.Site
	LinkService      *links.Service
	RatelimitService *ratelimit.Service
}

func CreateLink(d CreateLinkDeps) gin.HandlerFunc {
	qrHandler := CreateQR(CreateQRDeps{
		Log:              d.Log,
		Site:             d.Site,
		LinkService:      d.LinkService,
		RatelimitService: d.RatelimitService,
	})

	return func(c *gin.Context) {
		if c.Query("qr") != "" {
			qrHandler(c)
			return
		}

		url := c.Query("url")
		if url == "" {
			// p.IndexPage.WithError(c, "Bad url")
			d.Site.CreateLink(c)
			return
		}

		link, err := d.LinkService.CreateLink(c, url)
		if err == links.ErrBadUrl {
			d.Site.CreateLink(c)
			return
		}
		if err != nil {
			log.Error().Err(err).Msg("error creating link")
			d.Site.InternalError(c)
			return
		}

		d.Site.LinkResult(c, link)
	}
}
