package image

import "context"

type MetadataRepo interface {
	SaveImageMetadata(ctx context.Context, meta ImageMetadataDTO) error
	GetImageMetadataById(ctx context.Context, id string) (*ImageMetadataExDTO, error)
	GetImageMetadataDuplicate(ctx context.Context, size int, hash string) (*ImageMetadataExDTO, error)
}
