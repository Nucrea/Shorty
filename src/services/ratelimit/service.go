package ratelimit

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrTemporaryBanned = errors.New("ip temporary banned")
	ErrTooManyRequests = errors.New("too many requests")
)

const (
	LimitWindow = 1 * time.Minute
	LimitAmount = 30

	BanWindow = 30 * time.Minute
	BanAmount = 120
)

func NewService(rdb *redis.Client, log *zerolog.Logger, tracer trace.Tracer) *Service {
	newLog := log.With().Str("service", "ratelimit").Logger()
	return &Service{
		log:     &newLog,
		storage: &storage{rdb, tracer},
	}
}

type Service struct {
	log     *zerolog.Logger
	storage *storage
}

func (s *Service) Check(ctx context.Context, ip string) error {
	banned, err := s.storage.IsBanned(ctx, ip)
	if err != nil {
		s.log.Error().Err(err).Msgf("checking banned with storage")
		return err
	}
	if banned {
		s.log.Info().Msgf("rejected banned %s", ip)
		return ErrTemporaryBanned
	}

	rate, err := s.storage.IncRate(ctx, ip, LimitWindow)
	if err != nil {
		s.log.Error().Err(err).Msgf("inc requests rate with storage")
		return err
	}
	if rate >= BanAmount {
		s.storage.SetBanned(ctx, ip, BanWindow)
		s.log.Info().Msgf("temporary banned %s", ip)
		return ErrTemporaryBanned
	}
	if rate >= LimitAmount {
		s.log.Info().Msgf("too many requests from %s", ip)
		return ErrTooManyRequests
	}

	return nil
}
