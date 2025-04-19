package server

import (
	"fmt"
	"net/url"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *server) linkUrlfromId(linkId string) string {
	return fmt.Sprintf("%s/l/%s", s.Url, linkId)
}

func (s *server) LinkResult(c *gin.Context) templ.Component {
	inputUrl := c.PostForm("url")
	if inputUrl == "" {
		log.Error().Msg("empty url")
		c.Redirect(302, "/link?err="+url.QueryEscape("empty url"))
		return nil
	}
	return nil

	// var userIdPtr *string
	// userId := c.GetString("userId")
	// if userId != "" {
	// 	userIdPtr = &userId
	// }

	// link, err := s.LinksService.Create(c, inputUrl, userIdPtr)
	// if err == links.ErrBadUrl {
	// 	log.Error().Msg("bad url")
	// 	c.Redirect(302, "/link?err="+url.QueryEscape(err.Error()))
	// 	return nil
	// }
	// if err != nil {
	// 	log.Error().Err(err).Msg("error creating link")
	// 	s.site.InternalError(c)
	// 	return nil
	// }

	// resultUrl := fmt.Sprintf("%s/l/%s", s.Url, link.Id)
	// qrBase64, err := common.NewQRBase64(resultUrl)
	// if err != nil {
	// 	log.Error().Err(err).Msg("error creating qr")
	// 	s.site.InternalError(c)
	// 	return nil
	// }

	// return s.site.LinkResult(c, resultUrl, qrBase64)
}
