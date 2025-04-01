package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"shorty/server/middleware"
	"shorty/server/pages"
	"shorty/src/common/logging"
	"shorty/src/common/metrics"
	"shorty/src/common/tracing"
	"shorty/src/services/files"
	"shorty/src/services/guard"
	"shorty/src/services/image"
	"shorty/src/services/links"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

//go:embed static/*
var staticFS embed.FS

type Opts struct {
	Url          string
	ApiKey       string
	Logger       logging.Logger
	Tracer       trace.Tracer
	Meter        metrics.Meter
	LinksService *links.Service
	GuardService *guard.Service
	ImageService *image.Service
	FileService  *files.Service
}

func New(opts Opts) *server {
	opts.Logger = opts.Logger.WithService("server")
	return &server{
		opts,
		&pages.Site{},
		&profiler{
			mutex:        &sync.Mutex{},
			statusMetric: opts.Meter.NewGauge("profile_enabled", "Status flag of profiling mode (1-enabled, 0-disabled)"),
		},
	}
}

type server struct {
	Opts
	pages    *pages.Site
	profiler *profiler
}

func (s *server) Run(ctx context.Context, port uint16) {
	staticDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("opening static files dir")
	}

	gin.SetMode(gin.ReleaseMode)

	server := gin.New()
	server.ContextWithFallback = true // allows getting values from gin ctx, needed for tracing

	server.NoRoute(s.pages.NotFound)
	server.Use(middleware.Recovery(s.pages.InternalError, s.Logger, true))
	server.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	StaticFS(server, "/static", http.FS(staticDir))
	server.StaticFileFS("favicon.ico", "/favicon.ico", http.FS(staticDir))

	server.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/link")
	})

	server.Use(middleware.Log(s.Logger))
	server.Use(middleware.Metrics(s.Meter))
	server.Use(tracing.NewMiddleware(s.Tracer))
	server.Use(middleware.Ratelimit(s.GuardService, s.pages))
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{s.Url},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	profGroup := server.Group("/profile")
	{
		profGroup.Use(func(c *gin.Context) {
			if c.GetHeader("Authorization") != s.ApiKey {
				c.AbortWithStatus(403)
			} else {
				c.Next()
			}
		})
		profGroup.POST("/start", s.ProfileStart)
		profGroup.POST("/stop", s.ProfileStop)
	}

	server.GET("/link", s.pages.LinkForm)
	server.POST("/link", s.LinkResult)
	server.GET("/l/:id", s.LinkResolve)

	server.GET("/image", s.ImageForm)
	server.POST("/image", s.ImageUpload)
	server.GET("/image/view/:id", s.ImageView)
	server.GET("/i/:type/:id", s.ImageResolve)

	server.GET("/file", s.FileForm)
	server.POST("/file", s.FileUpload)
	server.GET("/file/view/:id", s.FileView)
	server.GET("/file/download/:id", s.FileDownload)
	server.GET("/f/:id/:name", s.FileResolve)

	s.Logger.Info().Msgf("Started server on port %d", port)
	server.Run(fmt.Sprintf(":%d", port))
}
