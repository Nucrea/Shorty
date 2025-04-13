package server

import (
	"shorty/src/services/links"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *server) LinkResolve(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		s.site.NotFound(c)
		return
	}

	url, err := s.LinksService.GetByShortId(c, id)
	if err == links.ErrNoSuchLink || err == links.ErrBadShortId {
		s.site.NotFound(c)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("error getting shortlink")
		s.site.InternalError(c)
		return
	}
	if url == "" {
		s.site.NotFound(c)
		return
	}

	c.Redirect(302, url)
}
