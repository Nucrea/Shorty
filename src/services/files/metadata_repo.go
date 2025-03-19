package files

import (
	"context"
	"shorty/src/common/assets"

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
	db           *pgxpool.Pool
	tracer       trace.Tracer
	assetStorage *assets.Storage
}

func (m *metadataRepo) SaveFileMetadata(ctx context.Context, meta FileMetadataDTO) error {
	_, span := m.tracer.Start(ctx, "postgres::SaveFileMetadata")
	defer span.End()

	query := `INSERT INTO files (id, file_id, name) VALUES ($1, $2, $3) RETURNING id;`
	_, err := m.db.Exec(ctx, query, meta.Id, meta.FileId, meta.Name)
	return err
}

func (m *metadataRepo) GetFileMetadata(ctx context.Context, id string) (*FileMetadataExDTO, error) {
	_, span := m.tracer.Start(ctx, "postgres::GetFileMetadata")
	defer span.End()

	query := `SELECT f.file_id, f.name, a.size, a.hash 
		FROM files f
		JOIN assets a on a.id = f.file_id
		WHERE f.id = $1;`
	row := m.db.QueryRow(ctx, query, id)

	dto := &FileMetadataExDTO{Id: id}
	err := row.Scan(&dto.FileId, &dto.Name, &dto.Size, &dto.Hash)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return dto, nil
}
