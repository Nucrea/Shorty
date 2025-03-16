package files

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

func newInfoStorage(db *pgxpool.Pool, tracer trace.Tracer) *infoStorage {
	return &infoStorage{
		db:     db,
		tracer: tracer,
	}
}

type infoStorage struct {
	db     *pgxpool.Pool
	tracer trace.Tracer
}

func (i *infoStorage) SaveFileInfo(ctx context.Context, dto FileInfoDTO) (*FileInfoDTO, error) {
	_, span := i.tracer.Start(ctx, "postgres::SaveFileInfo")
	defer span.End()

	query := `INSERT INTO files (short_id, resource_id, size, name) VALUES ($1, $2, $3, $4) RETURNING id;`
	row := i.db.QueryRow(ctx, query, dto.ShortId, dto.ResourceId, dto.Size, dto.Name)

	result := &FileInfoDTO{
		ShortId:    dto.ShortId,
		Size:       dto.Size,
		Name:       dto.Name,
		ResourceId: dto.ResourceId,
	}
	err := row.Scan(&result.Id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *infoStorage) GetFileInfo(ctx context.Context, shortId string) (*FileInfoDTO, error) {
	_, span := i.tracer.Start(ctx, "postgres::GetImageInfoByShortId")
	defer span.End()

	query := `SELECT id, resource_id, size, name FROM files WHERE short_id = $1;`
	row := i.db.QueryRow(ctx, query, shortId)

	result := &FileInfoDTO{ShortId: shortId}
	err := row.Scan(&result.Id, &result.ResourceId, &result.Size, &result.Name)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}
