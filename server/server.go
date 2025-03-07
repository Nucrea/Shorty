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
	"shorty/services/ban"
	"shorty/services/links"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

//go:embed static/*
var staticFS embed.FS

type ServerOpts struct {
	Port         uint16
	BaseUrl      string
	Log          *zerolog.Logger
	LinksService *links.Service
	BanService   *ban.Service
}

func Run(opts ServerOpts) {
	indexPage := index.NewPage()
	resultPage := result.NewPage()
	errorPage := genericerror.NewPage()

	server := gin.New()

	server.Use(gin.Recovery())
	server.Use(gin.Logger())

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
			Log:         opts.Log,
			IndexPage:   indexPage,
			ResultPage:  resultPage,
			ErrorPage:   errorPage,
			LinkService: opts.LinksService,
			BanService:  opts.BanService,
		},
	))
	server.GET("/:id", handlers.NewLinkResolveH(
		handlers.ResolveHDeps{
			BaseUrl:     opts.BaseUrl,
			Log:         opts.Log,
			LinkService: opts.LinksService,
			ErrorPage:   errorPage,
		},
	))

	server.Run(fmt.Sprintf(":%d", opts.Port))
}
