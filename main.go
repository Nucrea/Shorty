package main

import (
	"context"
	"fmt"
	"shorty/server"
	"shorty/src/common/logger"
	"shorty/src/common/tracing"
	"shorty/src/services/links"
	"shorty/src/services/ratelimit"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	conf, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("error parsing environment variables: %w", err))
	}

	log, err := logger.New(conf.LogFile)
	if err != nil {
		panic(err)
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

	tracer := tracing.NewNoopTracer()
	if conf.OTELUrl != "" {
		tracer, err = tracing.NewTracer(conf.OTELUrl)
		if err != nil {
			log.Fatal().Err(err).Msg("error init tracer")
		}
	}

	linksService := links.NewService(db, log, conf.AppUrl, tracer)
	ratelimitService := ratelimit.NewService(rdb, log, tracer)

	server.Run(server.ServerOpts{
		Port:             uint16(conf.AppPort),
		AppUrl:           conf.AppUrl,
		Log:              log,
		LinksService:     linksService,
		RatelimitService: ratelimitService,
		Tracer:           tracer,
	})
}
