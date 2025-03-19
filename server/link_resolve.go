package server

import (
	"shorty/src/services/links"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *server) LinkResolve(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		s.pages.NotFound(c)
		return
	}

	url, err := s.LinksService.GetByShortId(c, id)
	if err == links.ErrNoSuchLink || err == links.ErrBadShortId {
		s.pages.NotFound(c)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("error getting shortlink")
		s.pages.InternalError(c)
		return
	}
	if url == "" {
		s.pages.NotFound(c)
		return
	}

	c.Redirect(302, url)
}
