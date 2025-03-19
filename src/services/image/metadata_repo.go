package image

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

func newMetadataRepo(db *pgxpool.Pool, tracer trace.Tracer) *metadataRepo {
	return &metadataRepo{
		db:     db,
		tracer: tracer,
	}
}

type metadataRepo struct {
	db     *pgxpool.Pool
	tracer trace.Tracer
}

func (m *metadataRepo) SaveImageMetadata(ctx context.Context, meta ImageMetadataDTO) error {
	_, span := m.tracer.Start(ctx, "postgres::SaveImageInfo")
	defer span.End()

	query := `INSERT INTO images (id, name, original_id, thumbnail_id) VALUES ($1, $2, $3, $4);`
	_, err := m.db.Exec(ctx, query, meta.Id, meta.Name, meta.OriginalId, meta.ThumbnailId)
	return err
}

func (m *metadataRepo) GetImageMetadataById(ctx context.Context, id string) (*ImageMetadataExDTO, error) {
	_, span := m.tracer.Start(ctx, "postgres::GetImageInfoByShortId")
	defer span.End()

	query := `SELECT a.size, i.name, a.hash, i.original_id, i.thumbnail_id 
		FROM images i
		JOIN assets a ON a.id = i.original_id
		WHERE i.id = $1;`
	row := m.db.QueryRow(ctx, query, id)

	result := &ImageMetadataExDTO{Id: id}
	err := row.Scan(&result.Size, &result.Name, &result.Hash, &result.OriginalId, &result.ThumbnailId)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *metadataRepo) GetImageMetadataDuplicate(ctx context.Context, size int, hash string) (*ImageMetadataExDTO, error) {
	_, span := m.tracer.Start(ctx, "postgres::GetImageInfoByHash")
	defer span.End()

	query := `SELECT i.id, i.name, i.original_id, i.thumbnail_id 
		FROM images i
		JOIN assets a ON a.id = i.original_id
		WHERE a.hash = $1 AND a.size = $2;`
	row := m.db.QueryRow(ctx, query, hash, size)

	result := &ImageMetadataExDTO{Size: size, Hash: hash}
	err := row.Scan(&result.Id, &result.Name, &result.OriginalId, &result.ThumbnailId)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}
