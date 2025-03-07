package ratelimit

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
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

type Service struct {
	storage *storage
}

func NewService(rdb *redis.Client) *Service {
	return &Service{
		storage: &storage{rdb},
	}
}

func (s *Service) Check(ctx context.Context, ip string) error {
	banned, err := s.storage.IsBanned(ctx, ip)
	if err != nil {
		return err
	}
	if banned {
		return ErrTemporaryBanned
	}

	rate, err := s.storage.IncRate(ctx, ip, LimitWindow)
	if err != nil {
		return err
	}
	if rate >= BanAmount {
		s.storage.SetBanned(ctx, ip, BanWindow)
		return ErrTemporaryBanned
	}
	if rate >= LimitAmount {
		return ErrTooManyRequests
	}

	return nil
}
