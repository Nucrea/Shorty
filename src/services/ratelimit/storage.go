package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type storage struct {
	rdb *redis.Client
}

func NewStorage(rdb *redis.Client) *storage {
	return &storage{
		rdb: rdb,
	}
}

func (s *storage) IncRate(ctx context.Context, ip string, window time.Duration) (int, error) {
	key := fmt.Sprintf("rate:%s", ip)

	pipe := s.rdb.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.ExpireNX(ctx, key, window)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, err
	}

	return int(incr.Val()), nil
}

func (s *storage) IsBanned(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("banned:%s", ip)
	isBanned, err := s.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return isBanned > 0, nil
}

func (s *storage) SetBanned(ctx context.Context, ip string, banDuration time.Duration) error {
	key := fmt.Sprintf("banned:%s", ip)
	return s.rdb.Set(ctx, key, true, banDuration).Err()
}
