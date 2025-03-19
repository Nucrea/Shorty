package links

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

func NewStorage(conn *pgxpool.Pool, tracer trace.Tracer) *storage {
	return &storage{conn, tracer}
}

type storage struct {
	conn   *pgxpool.Pool
	tracer trace.Tracer
}

func (s *storage) CreateLink(ctx context.Context, id, url string) error {
	_, span := s.tracer.Start(ctx, "postgres::CreateLink")
	defer span.End()

	query := `insert into shortlinks(id, url) values($1, $2);`
	_, err := s.conn.Exec(ctx, query, id, url)
	return err
}

func (s *storage) GetLink(ctx context.Context, id string) (string, error) {
	_, span := s.tracer.Start(ctx, "postgres::GetLink")
	defer span.End()

	query := `
	update shortlinks 
		set read_count=read_count+1
		where id=$1
		returning url;`
	row := s.conn.QueryRow(ctx, query, id)

	var url string
	err := row.Scan(&url)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return url, nil
}
