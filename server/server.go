package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"shorty/server/middleware"
	"shorty/server/site"
	"shorty/src/common/logging"
	"shorty/src/common/metrics"
	"shorty/src/common/tracing"
	"shorty/src/services/files"
	"shorty/src/services/guard"
	"shorty/src/services/image"
	"shorty/src/services/links"
	"shorty/src/services/users"
	"time"

	"github.com/a-h/templ"
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
	UserService  *users.Service
}

func New(opts Opts) *server {
	opts.Logger = opts.Logger.WithService("server")
	return &server{
		opts,
		&site.Site{},
	}
}

type server struct {
	Opts
	site *site.Site
}

func (s *server) Run(ctx context.Context, port uint16) {
	staticDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("opening static files dir")
	}

	gin.SetMode(gin.ReleaseMode)

	server := gin.New()
	server.ContextWithFallback = true // allows getting values from gin ctx, needed for tracing

	server.NoRoute(s.TemplWrapper(s.site.NotFound))
	server.Use(middleware.Recovery(s.TemplWrapper(s.site.InternalError), s.Logger, s.Meter, true))
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
	server.Use(middleware.Ratelimit(s.GuardService, s.site))
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{s.Url},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	server.Use(middleware.Authorization(s.UserService))

	server.GET("/link", s.TemplWrapper(s.site.LinkForm))
	server.POST("/link", s.TemplWrapper(s.LinkResult))
	server.GET("/l/:id", s.TemplWrapper(s.LinkResolve))

	server.GET("/image", s.TemplWrapper(s.ImageForm))
	server.POST("/image", s.TemplWrapper(s.ImageUpload))
	server.GET("/image/view/:id", s.TemplWrapper(s.ImageView))
	server.GET("/i/:type/:id", s.TemplWrapper(s.ImageResolve))

	server.GET("/file", s.TemplWrapper(s.FileForm))
	server.POST("/file", s.TemplWrapper(s.FileUpload))
	server.GET("/file/view/:id", s.TemplWrapper(s.FileView))
	server.GET("/file/download/:id", s.TemplWrapper(s.FileDownload))
	server.GET("/f/:id/:name", s.TemplWrapper(s.FileResolve))

	server.GET("/login", s.TemplWrapper(s.site.UserLogin))
	server.POST("/logout", s.TemplWrapper(s.UserLogout))
	server.GET("/register", s.TemplWrapper(s.site.UserRegister))
	server.POST("/user/login", s.TemplWrapper(s.UserLogin))
	server.POST("/user/create", s.TemplWrapper(s.UserRegister))

	server.GET("/account", s.TemplWrapper(s.UserAccount))
	server.GET("/account/settings", s.TemplWrapper(s.UserAccount))
	server.GET("/account/links", s.TemplWrapper(s.UserAccount))

	s.Logger.Info().Msgf("Started server on port %d", port)
	server.Run(fmt.Sprintf(":%d", port))
}

type TemplFunc func(c *gin.Context) templ.Component

func (s *server) TemplWrapper(f TemplFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()

		children := f(c)
		if children == nil {
			return
		}

		var (
			user *users.UserDTO
			err  error
		)
		if session := middleware.GetUserSession(c); session != nil {
			user, err = s.UserService.GetById(c, session.UserId)
			if err != nil {
				s.site.InternalError(c)
				return
			}
		}

		var account *site.AccountData
		if user != nil {
			account = &site.AccountData{Email: user.Email}
		}
		duration := time.Since(start)
		ctx := templ.WithChildren(c, children)

		site.Layout(site.PageData{
			Account:        account,
			RenderDuration: duration,
		}).Render(ctx, c.Writer)
	}
}
