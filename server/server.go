package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"shorty/server/middleware"
	"shorty/server/pages"
	"shorty/src/common/logger"
	"shorty/src/common/tracing"
	"shorty/src/services/files"
	"shorty/src/services/image"
	"shorty/src/services/links"
	"shorty/src/services/ratelimit"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

//go:embed static/*
var staticFS embed.FS

type Opts struct {
	Url              string
	Log              logger.Logger
	Tracer           trace.Tracer
	LinksService     *links.Service
	RatelimitService *ratelimit.Service
	ImageService     *image.Service
	FileService      *files.Service
}

func New(opts Opts) *server {
	return &server{opts, &pages.Site{}}
}

type server struct {
	Opts
	pages *pages.Site
}

func (s *server) Run(ctx context.Context, port uint16) {
	staticDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		s.Log.Fatal().Err(err).Msg("opening static files dir")
	}

	gin.SetMode(gin.ReleaseMode)

	server := gin.New()
	server.ContextWithFallback = true // allows getting values from gin ctx, needed for tracing

	server.NoRoute(s.pages.NotFound)
	server.Use(middleware.Recovery(s.pages.InternalError, s.Log, true))
	server.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	server.StaticFS("/static", http.FS(staticDir))

	server.Use(middleware.Log(s.Log))
	server.Use(tracing.NewMiddleware(s.Tracer))
	server.Use(middleware.Ratelimit(s.RatelimitService, s.pages))

	server.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/link")
	})

	server.GET("/link", s.pages.LinkForm)
	server.POST("/link", s.LinkResult)
	server.GET("/l/:id", s.LinkResolve)

	server.GET("/image", s.pages.ImageForm)
	server.POST("/image", s.ImageUpload)
	server.GET("/image/view/:id", s.ImageView)
	server.GET("/i/f/:id", s.ImageResolve)
	server.GET("/i/t/:id", s.ImageResolve)

	server.GET("/file", s.pages.FileForm)
	server.POST("/file", s.FileUpload)
	server.GET("/file/view/:id", s.FileView)
	server.GET("/f/:id/:name", s.FileResolve)

	s.Log.Info().Msgf("Started server on port %d", port)
	server.Run(fmt.Sprintf(":%d", port))
}
