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
	ErrTemporaryBanned = errors.New("temporary banned")
	ErrTooManyRequests = errors.New("too many requests")
	ErrInternal        = errors.New("internal error")
)

const (
	LimitWindow = 1 * time.Minute
	LimitAmount = 60

	BanWindow = 1 * time.Hour
	BanAmount = 2 * LimitAmount
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
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "ratelimit::Check")
	defer span.End()

	banned, err := s.storage.IsBanned(ctx, ip)
	if err != nil {
		log.Error().Err(err).Msgf("checking banned with storage")
		return ErrInternal
	}
	if banned {
		log.Info().Msgf("rejected banned %s", ip)
		return ErrTemporaryBanned
	}

	rate, err := s.storage.IncRate(ctx, ip, LimitWindow)
	if err != nil {
		log.Error().Err(err).Msgf("inc requests rate with storage")
		return ErrInternal
	}
	if rate >= BanAmount {
		if err := s.storage.SetBanned(ctx, ip, BanWindow); err != nil {
			log.Error().Err(err).Msgf("set banned with storage")
			return ErrInternal
		}

		log.Info().Msgf("temporary banned %s", ip)
		return ErrTemporaryBanned
	}
	if rate >= LimitAmount {
		log.Info().Msgf("too many requests from %s", ip)
		return ErrTooManyRequests
	}

	return nil
}
