package main

import (
	"context"
	"os"
	"shorty/server"
	"shorty/src/services/ban"
	"shorty/src/services/links"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()

	log := zerolog.New(os.Stdout)

	conf, err := NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error parsing environment variables")
	}

	db, err := pgx.Connect(ctx, conf.PostgresUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to postgres")
	}

	redisOpts, err := redis.ParseURL(conf.RedisUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("error parsing redis url")
	}
	rdb := redis.NewClient(redisOpts)

	linksService := links.NewService(db, conf.AppUrl)
	banService := ban.NewService(rdb)

	server.Run(server.ServerOpts{
		Port:         uint16(conf.AppPort),
		AppUrl:       conf.AppUrl,
		Log:          &log,
		LinksService: linksService,
		BanService:   banService,
	})
}
