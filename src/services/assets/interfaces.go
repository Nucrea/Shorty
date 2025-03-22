package assets

import (
	"context"
)

type MetadataRepo interface {
	SaveAssetsMetadata(ctx context.Context, metas ...AssetMetadataDTO) error
	GetAssetMetadata(ctx context.Context, id string) (*AssetMetadataDTO, error)
	SetAssetsStatus(ctx context.Context, status AssetStatus, ids ...string) error
}
