package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"shorty/server/handlers"
	"shorty/server/site"
	"shorty/src/common/logger"
	"shorty/src/common/tracing"
	"shorty/src/services/image"
	"shorty/src/services/links"
	"shorty/src/services/ratelimit"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

//go:embed static/*
var staticFS embed.FS

type ServerOpts struct {
	Port             uint16
	AppUrl           string
	Log              logger.Logger
	LinksService     *links.Service
	RatelimitService *ratelimit.Service
	ImageService     *image.Service
	Tracer           trace.Tracer
}

func Run(opts ServerOpts) {
	site := &site.Site{}

	gin.SetMode(gin.ReleaseMode)

	server := gin.New()
	server.ContextWithFallback = true // Use it to allow getting values from c.Request.Context(). CRITICAL FOR TRACING

	server.Use(gin.Recovery())
	server.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	staticDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	server.StaticFS("/static", http.FS(staticDir))

	server.Use(NewRequestLogM(opts.Log))
	server.Use(tracing.NewMiddleware(opts.Tracer))
	server.Use(NewRatelimitM(opts.RatelimitService, site))

	server.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, fmt.Sprintf("%s/link", opts.AppUrl))
	})
	server.GET("/link", site.CreateLink)
	server.GET("/link/create", handlers.CreateLink(
		handlers.CreateLinkDeps{
			Log:              opts.Log,
			Site:             site,
			LinkService:      opts.LinksService,
			RatelimitService: opts.RatelimitService,
		},
	))
	server.GET("/s/:id", handlers.ResolveLink(
		handlers.ResolveLinkDeps{
			Log:         opts.Log,
			Site:        site,
			LinkService: opts.LinksService,
		},
	))
	server.GET("/image", site.UploadImage)
	server.POST("/image/upload", handlers.UploadImage(
		handlers.UploadImageDeps{
			BaseUrl:      opts.AppUrl,
			Log:          opts.Log,
			Site:         site,
			ImageService: opts.ImageService,
		},
	))
	server.GET("/image/view/:id", handlers.ViewImage(
		handlers.ViewImageDeps{
			BaseUrl:      opts.AppUrl,
			Log:          opts.Log,
			Site:         site,
			ImageService: opts.ImageService,
		},
	))
	server.GET("/i/f/:id", handlers.ResolveImage(
		handlers.ResolveImageDeps{
			Log:          opts.Log,
			Site:         site,
			ImageService: opts.ImageService,
		},
	))
	server.GET("/i/t/:id", handlers.ResolveImage(
		handlers.ResolveImageDeps{
			Log:          opts.Log,
			Site:         site,
			ImageService: opts.ImageService,
		},
	))

	opts.Log.Info().Msgf("Started server on port %d", opts.Port)
	server.Run(fmt.Sprintf(":%d", opts.Port))
}
