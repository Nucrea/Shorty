package postgres

import (
	"context"
	"fmt"
	"shorty/src/common/logging"
	"shorty/src/common/metrics"
	"shorty/src/services/assets"
	"shorty/src/services/files"
	"shorty/src/services/image"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

var DefaultBuckets = []float64{1, 10, 50, 100, 200}

func NewPostgres(ctx context.Context, connUrl string, logger logging.Logger, tracer trace.Tracer, meter metrics.Meter) (*Postgres, error) {
	dbPool, err := pgxpool.New(ctx, connUrl)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		db:     dbPool,
		tracer: tracer,
		logger: logger,
		meter:  meter,
	}, nil
}

type Postgres struct {
	db     *pgxpool.Pool
	tracer trace.Tracer
	logger logging.Logger
	meter  metrics.Meter
}

func (p *Postgres) SaveImageMetadata(ctx context.Context, meta image.ImageMetadataDTO) error {
	query := `INSERT INTO images (id, name, original_id, thumbnail_id) VALUES ($1, $2, $3, $4);`
	return exec(ctx, p, "SaveImageMetadata", query,
		meta.Id, meta.Name, meta.OriginalId, meta.ThumbnailId)
}

func (p *Postgres) GetImageMetadataDuplicate(ctx context.Context, size int, hash string) (*image.ImageMetadataExDTO, error) {
	scanFunc := func(row pgx.Row) (*image.ImageMetadataExDTO, error) {
		r := &image.ImageMetadataExDTO{Size: size, Hash: hash}
		return r, row.Scan(&r.Name, &r.OriginalResourceId, &r.OriginalId, &r.ThumbnailId, &r.ThumbnailResourceId)
	}

	query := `SELECT i.id, i.name, ao.id, ao.resource_id, at.id, at.resource_id
		FROM images i
		JOIN assets ao ON ao.id = i.original_id
		JOIN assets at ON at.id = i.thumbnail_id
		WHERE ao.hash = $1 AND ao.size = $2;`
	return queryRow(ctx, p, "GetImageMetadataDuplicate", scanFunc, query, hash, size)
}

func (p *Postgres) GetImageMetadataById(ctx context.Context, id string) (*image.ImageMetadataExDTO, error) {
	scanFunc := func(row pgx.Row) (*image.ImageMetadataExDTO, error) {
		r := &image.ImageMetadataExDTO{Id: id}
		return r, row.Scan(&r.Size, &r.Name, &r.Hash, &r.OriginalResourceId, &r.OriginalId, &r.ThumbnailId, &r.ThumbnailResourceId)
	}

	query := `SELECT ao.size, i.name, ao.hash, ao.id, ao.resource_id, at.id, at.resource_id 
		FROM images i
		JOIN assets ao ON ao.id = i.original_id
		JOIN assets at ON at.id = i.thumbnail_id
		WHERE i.id = $1;`
	return queryRow(ctx, p, "GetImageMetadataById", scanFunc, query, id)
}

func (p *Postgres) SaveShortlink(ctx context.Context, id, url string) error {
	query := `insert into shortlinks(id, url) values($1, $2);`
	return exec(ctx, p, "SaveShortlink", query, id, url)
}

func (p *Postgres) GetShortlink(ctx context.Context, id string) (string, error) {
	scanFunc := func(row pgx.Row) (string, error) {
		url := ""
		return url, row.Scan(&url)
	}

	query := `update shortlinks set read_count=read_count+1 where id=$1 returning url;`
	return queryRow(ctx, p, "GetShortlink", scanFunc, query, id)
}

func (p *Postgres) SaveFileMetadata(ctx context.Context, meta files.FileMetadataDTO) error {
	query := `INSERT INTO files (id, file_id, name) VALUES ($1, $2, $3) RETURNING id;`
	return exec(ctx, p, "SaveFileMetadata", query, meta.Id, meta.FileId, meta.Name)
}

func (p *Postgres) GetFileMetadata(ctx context.Context, id string) (*files.FileMetadataExDTO, error) {
	scanFunc := func(row pgx.Row) (*files.FileMetadataExDTO, error) {
		dto := &files.FileMetadataExDTO{Id: id}
		return dto, row.Scan(&dto.FileId, &dto.Name, &dto.Size, &dto.Hash)
	}

	query := `SELECT f.file_id, f.name, a.size, a.hash 
		FROM files f
		JOIN assets a on a.id = f.file_id
		WHERE f.id = $1;`
	return queryRow(ctx, p, "GetFileMetadata", scanFunc, query, id)
}

func (p *Postgres) GetAssetMetadata(ctx context.Context, id string) (*assets.AssetMetadataDTO, error) {
	scanFunc := func(row pgx.Row) (*assets.AssetMetadataDTO, error) {
		dto := &assets.AssetMetadataDTO{Id: id}
		return dto, row.Scan(&dto.ResourceId, &dto.Size, &dto.Hash, &dto.Bucket)
	}

	query := `select resource_id, size, hash, bucket from assets where id = $1;`
	return queryRow(ctx, p, "GetAssetMetadata", scanFunc, query, id)
}

func (p *Postgres) SaveAssetsMetadata(ctx context.Context, metas ...assets.AssetMetadataDTO) error {
	defer observe(ctx, p, "SaveAssetsMetadata")()

	rows := [][]any{}
	for _, meta := range metas {
		rows = append(rows, []any{meta.Id, meta.ResourceId, meta.Size, meta.Hash, meta.Bucket})
	}

	copyCount, err := p.db.CopyFrom(ctx,
		pgx.Identifier{"assets"},
		[]string{"id", "resource_id", "size", "hash", "bucket"},
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

func (p *Postgres) SetAssetsStatus(ctx context.Context, status assets.AssetStatus, ids ...string) error {
	query := `update assets set status=$1 where id in ($1);`
	return exec(ctx, p, "SetAssetsStatus", query, strings.Join(ids, ","))
}
