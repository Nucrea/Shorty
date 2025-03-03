package main

import (
	"context"
	"html/template"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type IndexPage struct {
	Error string
}

type ShortlinkPage struct {
	Shortlink string
}

func main() {
	ctx := context.Background()

	log := zerolog.New(os.Stdout)

	pgUrl := os.Getenv("SHORTY_POSTGRES_URL")
	if pgUrl == "" {
		log.Fatal().Msg("empty db url")
	}

	baseUrl := os.Getenv("SHORTY_BASE_URL")
	if baseUrl == "" {
		log.Fatal().Msg("empty base url")
	}

	db, err := pgx.Connect(ctx, pgUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to postgres")
	}

	storage := NewStorage(db)

	indexPageTmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		log.Fatal().Err(err).Msg("err parsing index page template")
	}

	linkPageTmpl, err := template.ParseFiles("views/shortlink.html")
	if err != nil {
		log.Fatal().Err(err).Msg("err parsing shortlink page template")
	}

	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(gin.Logger())

	server.GET("/", func(ctx *gin.Context) {
		indexPageTmpl.Execute(ctx.Writer, IndexPage{})
		ctx.Status(200)
		ctx.Header("Content-Type", "text/html")
	})

	server.GET("/create", func(ctx *gin.Context) {
		if url := ctx.Query("url"); url == "" {
			indexPageTmpl.Execute(ctx.Writer, IndexPage{Error: "Bad url"})
		} else {
			shortlinkId, err := storage.Create(ctx, url)
			if err != nil {
				log.Error().Err(err).Msg("error creating shortlink")
				indexPageTmpl.Execute(ctx.Writer, IndexPage{Error: "Error occured"})
			} else {
				link := baseUrl + "/" + shortlinkId
				linkPageTmpl.Execute(ctx.Writer, ShortlinkPage{Shortlink: link})
			}
		}

		ctx.Status(200)
		ctx.Header("Content-Type", "text/html")
	})

	server.GET("/:id", func(ctx *gin.Context) {
		shortId := ctx.Param("id")
		if shortId == "" {
			ctx.Status(404)
			return
		}

		url, err := storage.Get(ctx, shortId)
		if err != nil {
			log.Error().Err(err).Msg("error getting shortlink")
			ctx.Status(500)
			return
		}
		if url == "" {
			ctx.Status(404)
			return
		}

		ctx.Redirect(302, url)
	})

	server.Run(":8081")
}
