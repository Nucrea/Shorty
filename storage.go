package main

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

func (s *Storage) Create(ctx context.Context, url string) (string, error) {
	shortId := GenerateShortId(10)
	query := `insert into shortlinks(short_id, url) values($1, $2);`
	_, err := s.conn.Exec(ctx, query, shortId, url)
	return shortId, err
}

func (s *Storage) Get(ctx context.Context, shortId string) (string, error) {
	query := `select url from shortlinks where short_id=$1;`
	row := s.conn.QueryRow(ctx, query, shortId)

	var url string
	if err := row.Scan(&url); err != nil {
		return "", err
	}

	return url, nil
}
