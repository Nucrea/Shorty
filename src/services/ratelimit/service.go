package ratelimit

import (
	"context"
	"errors"
	"shorty/src/common/logger"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrTemporaryBanned = errors.New("ip temporary banned")
	ErrTooManyRequests = errors.New("too many requests")
	ErrInternal        = errors.New("internal error")
)

const (
	LimitWindow = 1 * time.Minute
	LimitAmount = 20

	BanWindow = 1 * time.Hour
	BanAmount = 60
)

func NewService(rdb *redis.Client, log logger.Logger, tracer trace.Tracer) *Service {
	return &Service{
		log:     log.WithService("ratelimit"),
		tracer:  tracer,
		storage: NewStorage(rdb, tracer),
	}
}

type Service struct {
	log     logger.Logger
	tracer  trace.Tracer
	storage *storage
}

func (s *Service) Check(ctx context.Context, ip string) error {
	ctx, span := s.tracer.Start(ctx, "ratelimit::Check")
	defer span.End()

	banned, err := s.storage.IsBanned(ctx, ip)
	if err != nil {
		s.log.Error().Err(err).Msgf("checking banned with storage")
		return ErrInternal
	}
	if banned {
		s.log.Info().Msgf("rejected banned %s", ip)
		return ErrTemporaryBanned
	}

	rate, err := s.storage.IncRate(ctx, ip, LimitWindow)
	if err != nil {
		s.log.Error().Err(err).Msgf("inc requests rate with storage")
		return ErrInternal
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
