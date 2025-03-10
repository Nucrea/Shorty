package image

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

type ImageInfoDTO struct {
	Id          string
	ShortId     string
	Size        int
	Name        string
	ImageId     string
	ThumbnailId string
}

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

	query := `INSERT INTO images (short_id, size, name, image_id, thumbnail_id) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	row := i.db.QueryRow(ctx, query, dto.ShortId, dto.Size, dto.Name, dto.ImageId, dto.ThumbnailId)

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

func (i *infoStorage) GetImageInfo(ctx context.Context, id string) (*ImageInfoDTO, error) {
	_, span := i.tracer.Start(ctx, "postgres::GetImageInfo")
	defer span.End()

	query := `SELECT id, short_id, size, name, image_id, thumbnail_id FROM images WHERE short_id = $1;`
	row := i.db.QueryRow(ctx, query, id)

	result := &ImageInfoDTO{}
	err := row.Scan(&result.Id, &result.ShortId, &result.Size, &result.Name, &result.ImageId, &result.ThumbnailId)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}
