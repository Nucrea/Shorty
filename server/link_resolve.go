package server

import (
	"shorty/src/services/links"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *server) LinkResolve(c *gin.Context) templ.Component {
	id := c.Param("id")
	if id == "" {
		return s.site.NotFound(c)
	}

	link, err := s.LinksService.GetById(c, id)
	if err == links.ErrNoSuchLink || err == links.ErrBadShortId {
		return s.site.NotFound(c)
	}
	if err != nil {
		log.Error().Err(err).Msg("error getting shortlink")
		return s.site.InternalError(c)
	}
	if link == nil {
		return s.site.NotFound(c)
	}

	c.Redirect(302, link.Url)
	return nil
}
