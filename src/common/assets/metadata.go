package assets

import (
	"context"
	"fmt"

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

func (m *metadataRepo) SaveAssetsMetadata(ctx context.Context, metas ...AssetMetadataDTO) error {
	_, span := m.tracer.Start(ctx, "postgres::SaveFileInfo")
	defer span.End()

	rows := [][]interface{}{}
	for _, meta := range metas {
		rows = append(rows, []interface{}{meta.Id, meta.Size, meta.Hash, meta.Bucket})
	}

	copyCount, err := m.db.CopyFrom(ctx,
		pgx.Identifier{"assets"},
		[]string{"id", "size", "hash", "bucket"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}
	if int(copyCount) != len(metas) {
		return fmt.Errorf("not all rows inserted")
	}

	return nil
}

func (m *metadataRepo) SetAssetsStatus(ctx context.Context, status AssetStatus, ids ...string) error {
	return nil
}
