package server

import (
	"shorty/server/middleware"
	"shorty/server/site/pages"
	"shorty/src/common"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func (s *server) UserAccount(c *gin.Context) templ.Component {
	session := middleware.GetUserSession(c)
	if session == nil {
		c.Redirect(302, "login")
		return nil
	}

	user, err := s.UserService.GetById(c, session.UserId)
	if err != nil {
		return s.site.InternalError(c)
	}
	links, _ := s.LinksService.GetByUserId(c, session.UserId)
	if err != nil {
		return s.site.InternalError(c)
	}

	newLinks := []pages.AccountViewLink{}
	for _, link := range links {
		sourceUrl := s.linkUrlfromId(link.Id)
		qrBase64, _ := common.NewQRBase64(sourceUrl)
		newLinks = append(newLinks, pages.AccountViewLink{
			SourceUrl:      sourceUrl,
			DestinationUrl: link.Url,
			QRBase64:       qrBase64,
		})
	}

	return s.site.AccountView(c, pages.AccountViewParams{
		User: pages.AccountViewUser{
			Email: user.Email,
		},
		Links: newLinks,
	})
}
