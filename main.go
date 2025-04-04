package main

import (
	"context"
	"fmt"
	"shorty/server"
	"shorty/src/common/logging"
	"shorty/src/common/metrics"
	"shorty/src/common/tracing"
	"shorty/src/databases/postgres"
	"shorty/src/databases/redis"
	"shorty/src/services/assets"
	"shorty/src/services/files"
	"shorty/src/services/guard"
	"shorty/src/services/image"
	"shorty/src/services/links"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx := context.Background()

	conf, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("error parsing environment variables: %w", err))
	}

	logger, err := logging.NewLogger(
		logging.WithFile(conf.LogFile),
		// logging.WithOpenTelemetry(conf.OTELUrl),
	)
	if err != nil {
		panic(err)
	}

	tracer, err := tracing.NewTracer(conf.OTELUrl)
	if err != nil {
		logger.Fatal().Err(err).Msg("error init tracer")
	}

	meter, err := metrics.NewOtel("shorty", conf.OTELUrl)
	if err != nil {
		logger.Fatal().Err(err).Msg("error init metrics")
	}

	pgdb, err := postgres.NewPostgres(ctx, conf.PostgresUrl, tracer, meter)
	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to postgres")
	}

	rdb, err := redis.New(ctx, conf.RedisUrl, tracer, meter)
	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to redis")
	}

	s3, err := minio.New(conf.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.MinioAccessKey, conf.MinioAccessSecret, ""),
		Secure: false,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("error init minio client")
	}

	assetsStorage := assets.NewStorage(pgdb, rdb, s3, tracer, logger)
	linksService := links.NewService(pgdb, logger, tracer, meter)
	guardService := guard.NewService(rdb, logger, tracer)
	imageService := image.NewService(pgdb, assetsStorage, logger, tracer)
	fileService := files.NewService(pgdb, assetsStorage, logger, tracer)

	srv := server.New(server.Opts{
		Url:          conf.AppUrl,
		ApiKey:       conf.ApiKey,
		Logger:       logger,
		Tracer:       tracer,
		Meter:        meter,
		LinksService: linksService,
		GuardService: guardService,
		ImageService: imageService,
		FileService:  fileService,
	})
	srv.Run(ctx, conf.AppPort)
}
