package files

import (
	"context"
	"errors"
	"shorty/src/common"
	"shorty/src/common/assets"
	"shorty/src/common/broker"
	"shorty/src/common/logging"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("file not found")
	ErrTooBig   = errors.New("file too big")
)

const (
	BucketName = "files"
	MaxSize    = 20 * 1024 * 1024
)

func NewService(pg *pgxpool.Pool, mc *minio.Client, log logging.Logger, tracer trace.Tracer) *Service {
	return &Service{
		log:          log.WithService("files"),
		tracer:       tracer,
		assetStorage: assets.NewStorage(pg, mc, tracer, log),
		infoStorage:  newMetadataRepo(pg, tracer),
	}
}

type Service struct {
	log          logging.Logger
	tracer       trace.Tracer
	broker       broker.Broker
	assetStorage *assets.Storage
	infoStorage  *metadataRepo
}

func (s *Service) UploadFile(ctx context.Context, name string, fileBytes []byte) (*FileMetadataDTO, error) {
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "files::UploadFile")
	defer span.End()

	if len(fileBytes) > MaxSize {
		return nil, ErrTooBig
	}

	result, err := s.assetStorage.SaveAssets(ctx, BucketName, fileBytes)
	if err != nil {
		log.Error().Err(err).Msg("err saving file asset")
		return nil, ErrInternal
	}

	metadata := &FileMetadataDTO{
		Id:     common.NewShortId(32),
		FileId: result[0].Id,
		Name:   name,
	}
	if err := s.infoStorage.SaveFileMetadata(ctx, *metadata); err != nil {
		log.Error().Err(err).Msg("err saving file info")
		return nil, ErrInternal
	}

	log.Info().Msgf("saved file with id=%s", metadata.Id)

	return metadata, nil
}

func (s *Service) GetFileMetadata(ctx context.Context, id string) (*FileMetadataExDTO, error) {
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "files::GetFileMetadata")
	defer span.End()

	meta, err := s.infoStorage.GetFileMetadata(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed getting file info")
		return nil, ErrInternal
	}
	if meta == nil {
		log.Info().Msgf("not found file with id=%s", id)
		return nil, ErrNotFound
	}

	log.Info().Msgf("read file metadata, id=%s", meta.Id)

	return meta, nil
}

func (s *Service) GetFileBytes(ctx context.Context, id string) ([]byte, error) {
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "files::GetFileBytes")
	defer span.End()

	meta, err := s.GetFileMetadata(ctx, id)
	if err != nil {
		return nil, err
	}

	assetBytes, err := s.assetStorage.GetAssetBytes(ctx, BucketName, meta.FileId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting file (id=%s, file_id=%s) from storage", id, meta.FileId)
		return nil, ErrInternal
	}

	log.Info().Msgf("read file bytes, id=%s", meta.Id)

	return assetBytes, nil
}
