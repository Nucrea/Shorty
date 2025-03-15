package server

import (
	"fmt"
	"shorty/src/services/links"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *server) LinkResult(c *gin.Context) {
	url := c.PostForm("url")
	if url == "" {
		// p.IndexPage.WithError(c, "Bad url")
		log.Error().Msg("empty url")
		s.pages.LinkForm(c)
		return
	}

	shortId, err := s.LinksService.Create(c, url)
	if err == links.ErrBadUrl {
		log.Error().Msg("bad url")
		s.pages.LinkForm(c)
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
