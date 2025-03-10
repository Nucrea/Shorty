package server

import (
	"embed"
	"encoding/base64"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"shorty/server/handlers"
	genericerror "shorty/server/pages/generic_error"
	"shorty/server/pages/index"
	"shorty/server/pages/result"
	"shorty/server/pages/upload"
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
	indexPage := index.NewPage()
	resultPage := result.NewPage()
	errorPage := genericerror.NewPage()
	uploadPage := upload.NewPage()

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
	server.Use(NewRatelimitM(opts.RatelimitService, errorPage))

	server.GET("/", indexPage.Clean)
	server.GET("/create", handlers.NewLinkCreateH(
		handlers.CreateHDeps{
			Log:              opts.Log,
			IndexPage:        indexPage,
			ResultPage:       resultPage,
			ErrorPage:        errorPage,
			LinkService:      opts.LinksService,
			RatelimitService: opts.RatelimitService,
			Tracer:           opts.Tracer,
		},
	))
	server.GET("/s/:id", handlers.NewLinkResolveH(
		handlers.ResolveHDeps{
			BaseUrl:     opts.AppUrl,
			Log:         opts.Log,
			LinkService: opts.LinksService,
			ErrorPage:   errorPage,
		},
	))
	server.GET("/image", uploadPage.Clean)
	server.POST("/image/create", handlers.NewImageCreateH(
		handlers.ImageHDeps{
			Log:          opts.Log,
			ResultPage:   resultPage,
			ErrorPage:    errorPage,
			ImageService: opts.ImageService,
			UploadPage:   uploadPage,
		},
	))
	server.GET("/image/:id", func(ctx *gin.Context) {
		bytes, err := opts.ImageService.GetThumbnail(ctx, ctx.Param("id"))
		if err != nil {
			errorPage.InternalError(ctx)
			return
		}

		aaa := base64.StdEncoding.EncodeToString(bytes)
		resultPage.WithQRCode(ctx, aaa)
	})

	opts.Log.Info().Msgf("Started server on port %d", opts.Port)
	server.Run(fmt.Sprintf(":%d", opts.Port))
}
