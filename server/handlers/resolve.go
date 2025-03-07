package handlers

import (
	genericerror "shorty/server/pages/generic_error"
	"shorty/src/services/links"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ResolveHDeps struct {
	BaseUrl     string
	Log         *zerolog.Logger
	ErrorPage   *genericerror.Page
	LinkService *links.Service
}

func NewLinkResolveH(d ResolveHDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortId := c.Param("id")
		if shortId == "" {
			c.Redirect(302, d.BaseUrl)
			return
		}

		url, err := d.LinkService.GetByShortId(c, shortId)
		if err == links.ErrNoSuchLink || err == links.ErrBadShortId {
			d.ErrorPage.NotFound(c)
			return
		}
		if err != nil {
			log.Error().Err(err).Msg("error getting shortlink")
			d.ErrorPage.InternalError(c)
			return
		}
		if url == "" {
			d.ErrorPage.NotFound(c)
			return
		}

		c.Redirect(302, url)
	}
}
