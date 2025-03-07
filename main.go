package main

import (
	"context"
	"os"
	"shorty/server"
	"shorty/src/services/links"
	"shorty/src/services/ratelimit"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()

	file, err := os.OpenFile("./.run/shorty.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		panic("error open log file")
	}
	defer file.Close()

	writer := zerolog.MultiLevelWriter(os.Stdout, file)

	log := zerolog.New(writer).With().Timestamp().Logger()

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

	tracer, err := NewTracer("http://localhost:4318")
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
