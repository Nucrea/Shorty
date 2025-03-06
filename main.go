package main

import (
	"context"
	"os"
	"shorty/handlers"
	genericerror "shorty/pages/generic_error"
	"shorty/pages/index"
	"shorty/pages/result"
	"shorty/services/ban"
	"shorty/services/links"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()

	log := zerolog.New(os.Stdout)

	pgUrl := os.Getenv("SHORTY_POSTGRES_URL")
	if pgUrl == "" {
		log.Fatal().Msg("empty db url")
	}

	redisUrl := os.Getenv("SHORTY_REDIS_URL")
	if redisUrl == "" {
		log.Fatal().Msg("empty redis url")
	}

	baseUrl := os.Getenv("SHORTY_BASE_URL")
	if baseUrl == "" {
		log.Fatal().Msg("empty base url")
	}

	db, err := pgx.Connect(ctx, pgUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to postgres")
	}

	redisOpts, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("error parsing redis url")
	}
	rdb := redis.NewClient(redisOpts)

	linkService := links.NewService(db, baseUrl)
	banService := ban.NewService(rdb)

	indexPage := index.NewPage()
	resultPage := result.NewPage()
	errorPage := genericerror.NewPage()

	server := gin.New()

	server.Use(gin.Recovery())
	server.Use(gin.Logger())

	server.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	server.Static("/static", "./static")

	server.GET("/", indexPage.Clean)
	server.GET("/create", handlers.NewLinkCreateH(
		handlers.CreateHDeps{
			Log:         &log,
			IndexPage:   indexPage,
			ResultPage:  resultPage,
			ErrorPage:   errorPage,
			LinkService: linkService,
			BanService:  banService,
		},
	))
	server.GET("/:id", handlers.NewLinkResolveH(
		handlers.ResolveHDeps{
			BaseUrl:     baseUrl,
			Log:         &log,
			LinkService: linkService,
			ErrorPage:   errorPage,
		},
	))

	server.Run(":8081")
}
