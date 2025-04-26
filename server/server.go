package server

import (
	"context"
	"fmt"
	"shorty/server/middleware"
	"shorty/src/common/logging"
	"shorty/src/common/metrics"
	"shorty/src/common/tracing"
	"shorty/src/services/files"
	"shorty/src/services/guard"
	"shorty/src/services/image"
	"shorty/src/services/links"
	"shorty/src/services/users"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// //go:embed static/*
// var staticFS embed.FS

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
	}
}

type server struct {
	Opts
}

func (s *server) Run(ctx context.Context, port uint16) {
	// staticDir, err := fs.Sub(staticFS, "static")
	// if err != nil {
	// 	s.Logger.Fatal().Err(err).Msg("opening static files dir")
	// }

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.ContextWithFallback = true // allows getting values from gin ctx, needed for tracing

	router.Use(middleware.Recovery(func(ctx *gin.Context) {
		ctx.Status(500)
	}, s.Logger, s.Meter, true))
	router.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	// StaticFS(server, "/static", http.FS(staticDir))
	// server.StaticFileFS("favicon.ico", "/favicon.ico", http.FS(staticDir))
	// server.GET("/", func(ctx *gin.Context) {
	// 	ctx.Redirect(302, "/link")
	// })

	router.Use(middleware.Log(s.Logger))
	router.Use(middleware.Metrics(s.Meter))
	router.Use(tracing.NewMiddleware(s.Tracer))
	router.Use(middleware.Ratelimit(s.GuardService))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{s.Url},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	// router.Use(middleware.Authorization(s.UserService))

	router.GET("/i/:type/:id", s.ResolveImage)

	api := router.Group("/api/v1")
	{
		linkRouter := api.Group("/link")
		linkRouter.POST("/create", wrap(s.Logger, s.PostLink))
		linkRouter.GET("/:id", wrap(s.Logger, s.GetLink))
	}
	{
		imageRouter := api.Group("/image")
		imageRouter.POST("/upload", wrap(s.Logger, s.UploadImage))
		imageRouter.GET("/info/:id", wrap(s.Logger, s.GetImageInfo))
	}
	{
		// fileRouter := api.Group("/file")
		// fileRouter.POST("/upload", func(ctx *gin.Context) {})
		// fileRouter.GET("/info/:id", func(ctx *gin.Context) {})
		// fileRouter.GET("/download/:id", func(ctx *gin.Context) {})
	}
	{
		userRouter := api.Group("/user")
		userRouter.POST("/login", wrap(s.Logger, s.LoginUser))
		userRouter.GET("/register", wrap(s.Logger, s.RegisterUser))
		userRouter.GET("/logout", wrap(s.Logger, s.LogoutUser))
	}

	s.Logger.Info().Msgf("Started server on port %d", port)
	router.Run(fmt.Sprintf(":%d", port))
}
