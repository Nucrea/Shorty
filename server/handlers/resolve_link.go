package handlers

import (
	"shorty/server/site"
	"shorty/src/common/logger"
	"shorty/src/services/links"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ResolveLinkDeps struct {
	Log         logger.Logger
	Site        *site.Site
	LinkService *links.Service
}

func ResolveLink(d ResolveLinkDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortId := c.Param("id")
		if shortId == "" {
			d.Site.NotFound(c)
			return
		}

		url, err := d.LinkService.GetByShortId(c, shortId)
		if err == links.ErrNoSuchLink || err == links.ErrBadShortId {
			d.Site.NotFound(c)
			return
		}
		if err != nil {
			log.Error().Err(err).Msg("error getting shortlink")
			d.Site.InternalError(c)
			return
		}
		if url == "" {
			d.Site.NotFound(c)
			return
		}

		c.Redirect(302, url)
	}
}
