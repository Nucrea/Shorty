package main

import (
	"context"
	"math"
	"os"
	"shorty/server"
	"shorty/services/ban"
	"shorty/services/links"
	"strconv"

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

	appPortEnv := os.Getenv("SHORTY_APP_PORT")
	if appPortEnv == "" {
		log.Fatal().Msg("empty app port")
	}
	appPort, err := strconv.Atoi(appPortEnv)
	if err != nil {
		log.Fatal().Err(err).Msg("error parsing app port")
	}
	if appPort > math.MaxUint16 {
		log.Fatal().Msg("app port is out of range")
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

	linksService := links.NewService(db, baseUrl)
	banService := ban.NewService(rdb)

	server.Run(server.ServerOpts{
		Port:         uint16(appPort),
		BaseUrl:      baseUrl,
		Log:          &log,
		LinksService: linksService,
		BanService:   banService,
	})
}
