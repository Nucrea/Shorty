package files

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"shorty/src/common"
	"shorty/src/common/assets"
	"shorty/src/common/broker"
	"shorty/src/common/logger"

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
	MaxSize = 20 * 1024 * 1024
)

func NewService(pg *pgxpool.Pool, mc *minio.Client, log logger.Logger, tracer trace.Tracer) *Service {
	return &Service{
		log:          log.WithService("files"),
		tracer:       tracer,
		assetStorage: assets.NewStorage(pg, mc, tracer, "files"),
		infoStorage:  newMetadataRepo(pg, tracer),
	}
}

type Service struct {
	log          logger.Logger
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

	hashBytes := sha512.Sum512(fileBytes)
	hash := hex.EncodeToString(hashBytes[:])

	asset := assets.AssetDTO{
		Id:    common.NewShortId(32),
		Size:  len(fileBytes),
		Hash:  hash,
		Bytes: fileBytes,
	}

	if _, err := s.assetStorage.SaveAssets(ctx, asset); err != nil {
		log.Error().Err(err).Msg("err saving file asset")
		return nil, ErrInternal
	}

	metadata := &FileMetadataDTO{
		Id:     common.NewShortId(32),
		FileId: asset.Id,
		Name:   name,
	}
	if err := s.infoStorage.SaveFileMetadata(ctx, *metadata); err != nil {
		log.Error().Err(err).Msg("err saving file info")
		return nil, ErrInternal
	}

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

	assetBytes, err := s.assetStorage.GetAssetBytes(ctx, meta.FileId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting file (id=%s, file_id=%s) from storage", id, meta.FileId)
		return nil, ErrInternal
	}

	return assetBytes, nil
}
