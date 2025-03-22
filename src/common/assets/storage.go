package assets

import (
	"context"
	"shorty/src/common"
	"shorty/src/common/logging"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
)

func NewStorage(pgxPool *pgxpool.Pool, s3 *minio.Client, tracer trace.Tracer, logger logging.Logger) *Storage {
	return &Storage{
		logger:   logger,
		tracer:   tracer,
		fileRepo: newFileRepo(s3, tracer),
		metaRepo: newMetadataRepo(pgxPool, tracer),
	}
}

type Storage struct {
	logger   logging.Logger
	tracer   trace.Tracer
	fileRepo *fileRepo
	metaRepo *metadataRepo
}

func (s *Storage) SaveAssets(ctx context.Context, bucket string, assets ...[]byte) ([]AssetMetadataDTO, error) {
	log := s.logger.WithContext(ctx).WithService("assets")

	ctx, span := s.tracer.Start(ctx, "assets::SaveAssets")
	defer span.End()

	ids := make([]string, len(assets))
	metadatas := make([]AssetMetadataDTO, len(assets))
	for i, asset := range assets {
		ids[i] = common.NewShortId(32)
		metadatas[i] = AssetMetadataDTO{
			Id:     ids[i],
			Size:   len(asset),
			Hash:   common.NewAssetHash(asset),
			Bucket: bucket,
		}
	}

	if err := s.metaRepo.SaveAssetsMetadata(ctx, metadatas...); err != nil {
		log.Error().Err(err).Msg("failed saving assets metadatas")
		return nil, err
	}

	for i, asset := range assets {
		meta := metadatas[i]
		if err := s.fileRepo.SaveFile(ctx, bucket, meta.Id, asset); err != nil {
			log.Error().Err(err).Msgf("failed saving asset file, bucket=%s", bucket)
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

	log.Info().Msgf("saved assets, bucket=%s, ids=%s", bucket, sb.String())

	return metadatas, nil
}

func (s *Storage) GetAssetBytes(ctx context.Context, bucket, id string) ([]byte, error) {
	log := s.logger.WithContext(ctx).WithService("assets")

	ctx, span := s.tracer.Start(ctx, "assets::GetAssetBytes")
	defer span.End()

	fileBytes, err := s.fileRepo.GetFile(ctx, bucket, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting asset bytes, bucket=%s, id=%s", bucket, id)
		return nil, err
	}

	log.Info().Msgf("got asset bytes, bucket=%s, id=%s", bucket, id)
	return fileBytes, nil
}

// func (s *Storage) GetAssetDuplicate(ctx context.Context, size int, hash string) (*AssetMetadataDTO, error) {
// 	ctx, span := s.tracer.Start(ctx, "assets::GetAsset")
// 	defer span.End()

// 	return nil, nil
// }
