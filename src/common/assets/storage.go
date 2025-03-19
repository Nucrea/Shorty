package assets

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
)

func NewStorage(pgxPool *pgxpool.Pool, s3 *minio.Client, tracer trace.Tracer, bucketName string) *Storage {
	return &Storage{
		pgxPool:  pgxPool,
		tracer:   tracer,
		fileRepo: newFileRepo(s3, tracer, bucketName),
		metaRepo: newMetadataRepo(pgxPool, tracer),
	}
}

type Storage struct {
	pgxPool  *pgxpool.Pool
	tracer   trace.Tracer
	fileRepo *fileRepo
	metaRepo *metadataRepo
}

func (s *Storage) SaveAssets(ctx context.Context, assets ...AssetDTO) (*AssetMetadataDTO, error) {
	ctx, span := s.tracer.Start(ctx, "assets::SaveAssets")
	defer span.End()

	metadatas := make([]AssetMetadataDTO, len(assets))
	for i, asset := range assets {
		metadatas[i] = AssetMetadataDTO{asset.Id, asset.Size, asset.Hash}
	}
	if err := s.metaRepo.SaveAssetsMetadata(ctx, metadatas...); err != nil {
		return nil, err
	}

	for _, a := range assets {
		if err := s.fileRepo.SaveFile(ctx, a.Id, a.Bytes); err != nil {
			return nil, err
		}
	}

	ids := make([]string, len(assets))
	for i, asset := range assets {
		ids[i] = asset.Id
	}
	if err := s.metaRepo.SetAssetsStatus(ctx, "created", ids...); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Storage) GetAssetBytes(ctx context.Context, id string) ([]byte, error) {
	ctx, span := s.tracer.Start(ctx, "assets::GetAssetBytes")
	defer span.End()

	return s.fileRepo.GetFile(ctx, id)
}

// func (s *Storage) GetAssetDuplicate(ctx context.Context, size int, hash string) (*AssetMetadataDTO, error) {
// 	ctx, span := s.tracer.Start(ctx, "assets::GetAsset")
// 	defer span.End()

// 	return nil, nil
// }
