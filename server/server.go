package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"shorty/server/handlers"
	genericerror "shorty/server/pages/generic_error"
	"shorty/server/pages/index"
	"shorty/server/pages/result"
	"shorty/src/services/links"
	"shorty/src/services/ratelimit"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

//go:embed static/*
var staticFS embed.FS

type ServerOpts struct {
	Port             uint16
	AppUrl           string
	Log              *zerolog.Logger
	LinksService     *links.Service
	RatelimitService *ratelimit.Service
}

func Run(opts ServerOpts) {
	indexPage := index.NewPage()
	resultPage := result.NewPage()
	errorPage := genericerror.NewPage()

	gin.SetMode(gin.ReleaseMode)
	server := gin.New()

	server.Use(gin.Recovery())
	server.Use(RequestLogM(opts.Log))

	server.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	staticDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	server.StaticFS("/static", http.FS(staticDir))

	server.GET("/", indexPage.Clean)
	server.GET("/create", handlers.NewLinkCreateH(
		handlers.CreateHDeps{
			Log:              opts.Log,
			IndexPage:        indexPage,
			ResultPage:       resultPage,
			ErrorPage:        errorPage,
			LinkService:      opts.LinksService,
			RatelimitService: opts.RatelimitService,
		},
	))
	server.GET("/:id", handlers.NewLinkResolveH(
		handlers.ResolveHDeps{
			BaseUrl:     opts.AppUrl,
			Log:         opts.Log,
			LinkService: opts.LinksService,
			ErrorPage:   errorPage,
		},
	))

	opts.Log.Info().Msgf("Started server on port %d", opts.Port)
	server.Run(fmt.Sprintf(":%d", opts.Port))
}
