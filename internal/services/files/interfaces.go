package files

import (
	"context"
)

type MetadataRepo interface {
	SaveFileMetadata(ctx context.Context, meta FileMetadataDTO) error
	GetFileMetadata(ctx context.Context, id string) (*FileMetadataExDTO, error)
}
