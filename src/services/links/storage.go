package links

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

type Storage struct {
	conn   *pgx.Conn
	tracer trace.Tracer
}

func NewStorage(conn *pgx.Conn, tracer trace.Tracer) *Storage {
	return &Storage{conn, tracer}
}

func (s *Storage) Create(ctx context.Context, shortId, url string) (string, error) {
	_, span := s.tracer.Start(ctx, "postgres::Create")
	defer span.End()

	query := `insert into shortlinks(short_id, url) values($1, $2);`
	_, err := s.conn.Exec(ctx, query, shortId, url)
	return shortId, err
}

func (s *Storage) Get(ctx context.Context, shortId string) (string, error) {
	_, span := s.tracer.Start(ctx, "postgres::Get")
	defer span.End()

	query := `
	update shortlinks 
		set read_count=read_count+1
		where short_id=$1
		returning url;`
	row := s.conn.QueryRow(ctx, query, shortId)

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
