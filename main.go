package main

import (
	"context"
	"fmt"
	"os"
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
	"shorty/src/services/users"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logging.Fatal(fmt.Errorf("panic occured on startup: %v", err))
			os.Exit(1)
		}
	}()

	ctx := context.Background()

	conf, err := NewConfig()
	if err != nil {
		logging.Fatal(fmt.Errorf("failed parsing environment variables: %w", err))
	}

	logger, err := logging.NewLogger(
		logging.WithFile(conf.LogFile),
		// logging.WithOpenTelemetry(conf.OTELUrl),
	)
	if err != nil {
		logging.Fatal(fmt.Errorf("failed initializing logger: %w", err))
	}

	tracer, err := tracing.NewTracer(conf.OTELUrl)
	if err != nil {
		logger.Fatal().Err(err).Msg("error init tracer")
	}

	meter, err := metrics.NewOtel("shorty", conf.OTELUrl)
	if err != nil {
		logger.Fatal().Err(err).Msg("error init metrics")
	}

	pgdb, err := postgres.NewPostgres(ctx, conf.PostgresUrl, logger, tracer, meter)
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

	assetsStorage := assets.NewStorage(pgdb, rdb, s3, logger, tracer)
	linksService := links.NewService(pgdb, logger, tracer, meter)
	guardService := guard.NewService(rdb, logger, tracer, meter)
	imageService := image.NewService(pgdb, assetsStorage, logger, tracer, meter)
	fileService := files.NewService(pgdb, assetsStorage, logger, tracer, meter)
	userService := users.NewService(pgdb, rdb, logger, tracer, meter)

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
		UserService:  userService,
	})
	srv.Run(ctx, conf.AppPort)
}
