package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"shorty/server"
	"shorty/src/common/tracing"
	"shorty/src/services/links"
	"shorty/src/services/ratelimit"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()

	conf, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("error parsing environment variables: %w", err))
	}

	var writer io.Writer = os.Stdout

	if conf.LogFile != "" {
		file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			panic("error opening log file")
		}
		defer file.Close()

		writer = io.MultiWriter(writer, file)
	}

	log := zerolog.New(writer).With().Timestamp().Logger()

	db, err := pgx.Connect(ctx, conf.PostgresUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to postgres")
	}

	redisOpts, err := redis.ParseURL(conf.RedisUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("error parsing redis url")
	}
	rdb := redis.NewClient(redisOpts)

	tracer, err := tracing.NewTracer(conf.OTELUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("error init tracer")
	}

	linksService := links.NewService(db, conf.AppUrl, tracer)
	ratelimitService := ratelimit.NewService(rdb)

	server.Run(server.ServerOpts{
		Port:             uint16(conf.AppPort),
		AppUrl:           conf.AppUrl,
		Log:              &log,
		LinksService:     linksService,
		RatelimitService: ratelimitService,
		Tracer:           tracer,
	})
}
