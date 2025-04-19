package postgres

import (
	"context"
	"fmt"
	"shorty/src/common/logging"
	"shorty/src/common/metrics"
	"shorty/src/services/assets"
	"shorty/src/services/files"
	"shorty/src/services/image"
	"shorty/src/services/links"
	"shorty/src/services/users"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

var _ users.UserRepo = (*Postgres)(nil)
var _ assets.MetadataRepo = (*Postgres)(nil)
var _ image.MetadataRepo = (*Postgres)(nil)
var _ files.MetadataRepo = (*Postgres)(nil)
var _ links.Storage = (*Postgres)(nil)

func NewPostgres(
	ctx context.Context, connUrl string,
	logger logging.Logger, tracer trace.Tracer, meter metrics.Meter,
) (*Postgres, error) {
	dbPool, err := pgxpool.New(ctx, connUrl)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		db:     dbPool,
		logger: logger,
		tracer: tracer,
		meter:  meter,
		latencyHist: meter.NewHistogram(
			"postgres_query_latency",
			"Latency of postgresql queries",
			[]float64{1, 10, 50, 100, 200},
		),
	}, nil
}

type Postgres struct {
	tracer trace.Tracer
	logger logging.Logger
	meter  metrics.Meter

	db          *pgxpool.Pool
	latencyHist metrics.Histogram
}

// DeleteLinks implements links.Storage.
func (p *Postgres) DeleteLinks(ctx context.Context, ids ...string) error {
	panic("unimplemented")
}

func (p *Postgres) CreateUser(ctx context.Context, user users.UserDTO) (*users.UserDTO, error) {
	scanFunc := func(row pgx.Row) (*users.UserDTO, error) {
		r := &users.UserDTO{Email: user.Email, Secret: user.Secret}
		return r, row.Scan(&r.Id)
	}

	query := `INSERT INTO users (email, secret) VALUES ($1, $2) RETURNING id;`
	return queryRow(ctx, p, "CreateUser", scanFunc, query, user.Email, user.Secret)
}

func (p *Postgres) GetUserByEmail(ctx context.Context, email string) (*users.UserDTO, error) {
	scanFunc := func(row pgx.Row) (*users.UserDTO, error) {
		r := &users.UserDTO{Email: email}
		return r, row.Scan(&r.Id, &r.Secret, &r.Verified)
	}

	query := `SELECT id, secret, verified FROM users WHERE email=$1;`
	return queryRow(ctx, p, "GetUserByEmail", scanFunc, query, email)
}

func (p *Postgres) GetUserById(ctx context.Context, id string) (*users.UserDTO, error) {
	scanFunc := func(row pgx.Row) (*users.UserDTO, error) {
		r := &users.UserDTO{Id: id}
		return r, row.Scan(&r.Email, &r.Secret, &r.Verified)
	}

	query := `SELECT email, secret, verified FROM users WHERE id=$1;`
	return queryRow(ctx, p, "GetUserById", scanFunc, query, id)
}

func (p *Postgres) SetUserVerified(ctx context.Context, id string) error {
	query := `UPDATE users SET verified=true WHERE id=$1;`
	return exec(ctx, p, "SaveImageMetadata", query, id)
}

func (p *Postgres) SaveImageMetadata(ctx context.Context, meta image.ImageMetadataDTO) error {
	query := `INSERT INTO images (id, name, original_id, thumbnail_id) VALUES ($1, $2, $3, $4);`
	return exec(ctx, p, "SaveImageMetadata", query,
		meta.Id, meta.Name, meta.OriginalId, meta.ThumbnailId)
}

func (p *Postgres) GetImageMetadataDuplicate(ctx context.Context, size int, hash string) (*image.ImageMetadataExDTO, error) {
	scanFunc := func(row pgx.Row) (*image.ImageMetadataExDTO, error) {
		r := &image.ImageMetadataExDTO{Size: size, Hash: hash}
		return r, row.Scan(&r.Id, &r.Name, &r.OriginalId, &r.OriginalResourceId, &r.ThumbnailId, &r.ThumbnailResourceId)
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

func (p *Postgres) GetLinkById(ctx context.Context, id string) (*links.LinkDTO, error) {
	scanFunc := func(row pgx.Row) (*links.LinkDTO, error) {
		dto := &links.LinkDTO{Id: id}
		return dto, row.Scan(&dto.UserId, &dto.Url)
	}

	query := `select coalesce(user_id::text, ''), url from shortlinks where id=$1;`
	return queryRow(ctx, p, "GetLinkById", scanFunc, query, id)
}

func (p *Postgres) GetLinksByUserId(ctx context.Context, userId string) ([]*links.LinkDTO, error) {
	scanFunc := func(row pgx.Row) (*links.LinkDTO, error) {
		dto := &links.LinkDTO{UserId: userId}
		return dto, row.Scan(&dto.Id, &dto.Url)
	}

	query := `select id, url from shortlinks where user_id=$1;`
	return queryRows(ctx, p, "GetLinksByUserId", scanFunc, query, userId)
}

func (p *Postgres) SaveLink(ctx context.Context, id, url string) error {
	query := `insert into shortlinks(id, url) values($1, $2);`
	return exec(ctx, p, "SaveShortlink", query, id, url)
}

func (p *Postgres) SaveLinkForUser(ctx context.Context, id, userId, url string) error {
	query := `insert into shortlinks(id, user_id, url) values($1, $2, $3);`
	return exec(ctx, p, "SaveShortlink", query, id, userId, url)
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
	query := `update assets set status=$1 where id in ($2);`
	return exec(ctx, p, "SetAssetsStatus", query, status, strings.Join(ids, ","))
}
