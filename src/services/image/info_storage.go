package image

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

func newInfoStorage(db *pgx.Conn, tracer trace.Tracer) *infoStorage {
	return &infoStorage{
		db:     db,
		tracer: tracer,
	}
}

type infoStorage struct {
	db     *pgx.Conn
	tracer trace.Tracer
}

func (i *infoStorage) SaveImageInfo(ctx context.Context, dto ImageInfoDTO) (*ImageInfoDTO, error) {
	_, span := i.tracer.Start(ctx, "postgres::SaveImageInfo")
	defer span.End()

	query := `INSERT INTO images (short_id, size, name, hash, image_id, thumbnail_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
	row := i.db.QueryRow(ctx, query, dto.ShortId, dto.Size, dto.Name, dto.Hash, dto.ImageId, dto.ThumbnailId)

	result := &ImageInfoDTO{
		ShortId:     dto.ShortId,
		Size:        dto.Size,
		Name:        dto.Name,
		ImageId:     dto.ImageId,
		ThumbnailId: dto.ThumbnailId,
	}
	err := row.Scan(&result.Id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *infoStorage) GetImageInfoByShortId(ctx context.Context, shortId string) (*ImageInfoDTO, error) {
	_, span := i.tracer.Start(ctx, "postgres::GetImageInfoByShortId")
	defer span.End()

	query := `SELECT id, size, name, hash, image_id, thumbnail_id FROM images WHERE short_id = $1;`
	row := i.db.QueryRow(ctx, query, shortId)

	result := &ImageInfoDTO{ShortId: shortId}
	err := row.Scan(&result.Id, &result.Size, &result.Name, &result.Hash, &result.ImageId, &result.ThumbnailId)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *infoStorage) GetImageInfoByHash(ctx context.Context, hash string) (*ImageInfoDTO, error) {
	_, span := i.tracer.Start(ctx, "postgres::GetImageInfoByHash")
	defer span.End()

	query := `SELECT id, short_id, size, name, image_id, thumbnail_id FROM images WHERE hash = $1;`
	row := i.db.QueryRow(ctx, query, hash)

	result := &ImageInfoDTO{Hash: hash}
	err := row.Scan(&result.Id, &result.ShortId, &result.Size, &result.Name, &result.ImageId, &result.ThumbnailId)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}
