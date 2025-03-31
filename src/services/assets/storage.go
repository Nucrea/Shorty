package assets

import (
	"context"
	"shorty/src/common"
	"shorty/src/common/logging"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

func NewStorage(metaRepo MetadataRepo, s3 *minio.Client, rdb *redis.Client, tracer trace.Tracer, logger logging.Logger) *Storage {
	return &Storage{
		logger:    logger,
		tracer:    tracer,
		fileRepo:  newFileRepo(s3, tracer),
		metaRepo:  metaRepo,
		metaCache: &cache{rdb, tracer},
	}
}

type Storage struct {
	logger    logging.Logger
	tracer    trace.Tracer
	fileRepo  *fileRepo
	metaRepo  MetadataRepo
	metaCache *cache
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
			Id:         ids[i],
			ResourceId: common.NewShortId(32),
			Size:       len(asset),
			Hash:       common.NewAssetHash(asset),
			Bucket:     bucket,
		}
	}

	if err := s.metaRepo.SaveAssetsMetadata(ctx, metadatas...); err != nil {
		log.Error().Err(err).Msg("failed saving assets metadatas")
		return nil, err
	}

	for i, asset := range assets {
		meta := metadatas[i]
		if err := s.fileRepo.SaveFile(ctx, bucket, meta.ResourceId, asset); err != nil {
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

	for _, meta := range metadatas {
		if err := s.metaCache.PutMetadata(ctx, meta); err != nil {
			s.logger.Warning().Err(err).Msg("failed putting metadata to cache")
		}
	}

	return metadatas, nil
}

func (s *Storage) getAssetMetadata(ctx context.Context, id string) (*AssetMetadataDTO, error) {
	if meta, err := s.metaCache.GetMetadata(ctx, id); err == nil && meta != nil {
		return meta, nil
	}

	meta, err := s.metaRepo.GetAssetMetadata(ctx, id)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, nil
	}

	if err := s.metaCache.PutMetadata(ctx, *meta); err != nil {
		s.logger.Warning().Err(err).Msg("failed putting metadata to cache")
	}
	return meta, nil
}

func (s *Storage) GetAssetBytes(ctx context.Context, bucket, id string) ([]byte, error) {
	log := s.logger.WithContext(ctx).WithService("assets")

	ctx, span := s.tracer.Start(ctx, "assets::GetAssetBytes")
	defer span.End()

	meta, err := s.getAssetMetadata(ctx, id)
	if meta == nil {
		log.Info().Msgf("no such asset, bucket=%s id=%s", bucket, id)
	}

	fileBytes, err := s.fileRepo.GetFile(ctx, bucket, meta.ResourceId)
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
