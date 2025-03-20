package assets

import (
	"context"
	"shorty/src/common"
	"shorty/src/common/logger"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
)

func NewStorage(pgxPool *pgxpool.Pool, s3 *minio.Client, tracer trace.Tracer, logger logger.Logger, bucketName string) *Storage {
	return &Storage{
		logger:   logger,
		tracer:   tracer,
		fileRepo: newFileRepo(s3, tracer, bucketName),
		metaRepo: newMetadataRepo(pgxPool, tracer),
	}
}

type Storage struct {
	logger   logger.Logger
	tracer   trace.Tracer
	fileRepo *fileRepo
	metaRepo *metadataRepo
}

func (s *Storage) SaveAssets(ctx context.Context, assets ...[]byte) ([]AssetMetadataDTO, error) {
	log := s.logger.WithContext(ctx).WithService("assets")

	ctx, span := s.tracer.Start(ctx, "assets::SaveAssets")
	defer span.End()

	ids := make([]string, len(assets))
	metadatas := make([]AssetMetadataDTO, len(assets))
	for i, asset := range assets {
		ids[i] = common.NewShortId(32)
		metadatas[i] = AssetMetadataDTO{
			Id:   ids[i],
			Size: len(asset),
			Hash: common.NewAssetHash(asset),
		}
	}

	if err := s.metaRepo.SaveAssetsMetadata(ctx, metadatas...); err != nil {
		log.Error().Err(err).Msg("failed saving assets metadatas")
		return nil, err
	}

	for i, asset := range assets {
		meta := metadatas[i]
		if err := s.fileRepo.SaveFile(ctx, meta.Id, asset); err != nil {
			log.Error().Err(err).Msg("failed saving asset file")
			return nil, err
		}
	}

	if err := s.metaRepo.SetAssetsStatus(ctx, "created", ids...); err != nil {
		log.Error().Err(err).Msg("failed updating assets statuses")
		return nil, err
	}

	sb := strings.Builder{}
	sb.WriteRune('[')
	for _, id := range ids {
		sb.WriteString(id)
		sb.WriteString(", ")
	}
	sb.WriteRune(']')

	log.Info().Msgf("saved assets, ids=%s", sb.String())

	return metadatas, nil
}

func (s *Storage) GetAssetBytes(ctx context.Context, id string) ([]byte, error) {
	log := s.logger.WithContext(ctx).WithService("assets")

	ctx, span := s.tracer.Start(ctx, "assets::GetAssetBytes")
	defer span.End()

	fileBytes, err := s.fileRepo.GetFile(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting asset bytes, id=%s", id)
		return nil, err
	}

	log.Info().Msgf("got asset bytes, id=%s", id)
	return fileBytes, nil
}

// func (s *Storage) GetAssetDuplicate(ctx context.Context, size int, hash string) (*AssetMetadataDTO, error) {
// 	ctx, span := s.tracer.Start(ctx, "assets::GetAsset")
// 	defer span.End()

// 	return nil, nil
// }
