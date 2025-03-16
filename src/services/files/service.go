package files

import (
	"context"
	"errors"
	"shorty/src/common/broker"
	"shorty/src/common/logger"
	"shorty/src/common/s3"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("file not found")
	ErrTooBig   = errors.New("file too big")
)

func NewService(pg *pgxpool.Pool, mc *minio.Client, log logger.Logger, tracer trace.Tracer) *Service {
	return &Service{
		log:         log.WithService("files"),
		tracer:      tracer,
		fileStorage: s3.NewFileStorage(mc, tracer, "files"),
		infoStorage: newInfoStorage(pg, tracer),
	}
}

type Service struct {
	log         logger.Logger
	tracer      trace.Tracer
	broker      broker.Broker
	fileStorage *s3.FileStorage
	infoStorage *infoStorage
}

func (s *Service) UploadFile(ctx context.Context, name string, fileBytes []byte) (*FileInfoDTO, error) {
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "files::UploadFile")
	defer span.End()

	if len(fileBytes) > 20*1024*1024 {
		return nil, ErrTooBig
	}

	resourceId := NewShortId(32)
	if err := s.fileStorage.SaveFile(ctx, resourceId, fileBytes); err != nil {
		log.Error().Err(err).Msg("err saving file")
		return nil, ErrInternal
	}

	shortId := NewShortId(32)
	result, err := s.infoStorage.SaveFileInfo(ctx, FileInfoDTO{
		ShortId:    shortId,
		Name:       name,
		Size:       len(fileBytes),
		ResourceId: resourceId,
	})
	if err != nil {
		log.Error().Err(err).Msg("err saving file info")
		return nil, ErrInternal
	}

	return result, nil
}

func (s *Service) GetFileInfo(ctx context.Context, shortId string) (*FileInfoDTO, error) {
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "files::GetFileInfo")
	defer span.End()

	fileInfo, err := s.infoStorage.GetFileInfo(ctx, shortId)
	if err != nil {
		log.Error().Err(err).Msg("failed getting file info")
		return nil, ErrInternal
	}
	if fileInfo == nil {
		log.Info().Msgf("not found file with id=%s", shortId)
		return nil, ErrNotFound
	}

	return fileInfo, nil
}

func (s *Service) GetFile(ctx context.Context, resourceId string) ([]byte, error) {
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "files::GetFileInfo")
	defer span.End()

	fileBytes, err := s.fileStorage.GetFile(ctx, resourceId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting file (id=%s) from storage", resourceId)
		return nil, ErrInternal
	}

	return fileBytes, nil
}
