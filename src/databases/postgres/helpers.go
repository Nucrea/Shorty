package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func observe(ctx context.Context, p *Postgres, funcName string) func() {
	hist := p.meter.NewHistogram(funcName, "test", DefaultBuckets)

	_, span := p.tracer.Start(ctx, fmt.Sprintf("postgres::%s", funcName))
	defer span.End()

	start := time.Now()
	return func() {
		span.End()
		duration := time.Now().Sub(start).Milliseconds()
		hist.Observe(float64(duration))
	}
}

func queryRow[T any](
	ctx context.Context,
	p *Postgres,
	funcName string,
	scanFunc func(row pgx.Row) (T, error),
	query string, arguments ...any,
) (T, error) {
	defer observe(ctx, p, funcName)()

	val, err := scanFunc(p.db.QueryRow(ctx, query, arguments...))

	var t T
	if err == pgx.ErrNoRows {
		return t, nil
	}
	if err != nil {
		p.logger.WithContext(ctx).Error().Err(err).Str("func", funcName).Msg("failed exec db query")
		return t, err
	}

	return val, nil
}

func exec(ctx context.Context, p *Postgres, funcName, query string, arguments ...any) error {
	defer observe(ctx, p, funcName)()

	_, err := p.db.Exec(ctx, query, arguments...)
	if err != nil {
		p.logger.WithContext(ctx).Error().Err(err).Str("func", funcName).Msg("failed exec db query")
	}
	return err
}
