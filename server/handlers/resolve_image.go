package handlers

import (
	"shorty/server/site"
	"shorty/src/common/logger"
	"shorty/src/services/image"

	"github.com/gin-gonic/gin"
)

type ResolveImageDeps struct {
	Log          logger.Logger
	Site         *site.Site
	ImageService *image.Service
}

func ResolveImage(d ResolveImageDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortId := c.Param("id")
		if shortId == "" {
			c.AbortWithStatus(404)
			return
		}

		img, err := d.ImageService.GetFile(c, shortId)
		if err == image.ErrImageNotFound {
			c.AbortWithStatus(404)
			return
		}
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		c.Data(200, "image/jpeg", img)
	}
}
