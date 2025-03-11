package handlers

import (
	"fmt"
	"shorty/server/site"
	"shorty/src/common/logger"
	"shorty/src/services/image"

	"github.com/gin-gonic/gin"
)

type ViewImageDeps struct {
	BaseUrl      string
	Log          logger.Logger
	Site         *site.Site
	ImageService *image.Service
}

func ViewImage(d ViewImageDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortId := c.Param("id")
		if shortId == "" {
			d.Site.NotFound(c)
			return
		}

		info, err := d.ImageService.GetImageInfo(c, shortId)
		if err == image.ErrImageNotFound {
			d.Site.NotFound(c)
			return
		}
		if err != nil {
			d.Site.InternalError(c)
			return
		}

		imgUrl := fmt.Sprintf("%s/i/f/%s", d.BaseUrl, info.ImageId)
		d.Site.ViewImage(c, imgUrl, info.Name, info.Size)
	}
}
