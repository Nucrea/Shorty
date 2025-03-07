package links

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(conn *pgx.Conn) *Storage {
	return &Storage{conn}
}

func (s *Storage) Create(ctx context.Context, shortId, url string) (string, error) {
	query := `insert into shortlinks(short_id, url) values($1, $2);`
	_, err := s.conn.Exec(ctx, query, shortId, url)
	return shortId, err
}

func (s *Storage) Get(ctx context.Context, shortId string) (string, error) {
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
