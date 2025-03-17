package server

import (
	"fmt"
	"net/url"
	"shorty/src/services/links"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *server) LinkResult(c *gin.Context) {
	inputUrl := c.PostForm("url")
	if inputUrl == "" {
		log.Error().Msg("empty url")
		c.Redirect(302, "/link?err="+url.QueryEscape("empty url"))
		return
	}

	shortId, err := s.LinksService.Create(c, inputUrl)
	if err == links.ErrBadUrl {
		log.Error().Msg("bad url")
		c.Redirect(302, "/link?err="+url.QueryEscape(err.Error()))
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("error creating link")
		s.pages.InternalError(c)
		return
	}

	resultUrl := fmt.Sprintf("%s/l/%s", s.Url, shortId)
	qrBase64, err := s.LinksService.MakeQR(c, resultUrl)
	if err != nil {
		log.Error().Err(err).Msg("error creating qr")
		s.pages.InternalError(c)
		return
	}

	s.pages.LinkResult(c, resultUrl, qrBase64)
}
