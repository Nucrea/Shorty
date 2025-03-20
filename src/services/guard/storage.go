package guard

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

type storage struct {
	rdb    *redis.Client
	tracer trace.Tracer
}

func newStorage(rdb *redis.Client, tracer trace.Tracer) *storage {
	return &storage{
		rdb:    rdb,
		tracer: tracer,
	}
}

func (s *storage) IncRate(ctx context.Context, ip string, window time.Duration) (int, error) {
	ctx, span := s.tracer.Start(ctx, "redis::IncRate")
	defer span.End()

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
	ctx, span := s.tracer.Start(ctx, "redis::IsBanned")
	defer span.End()

	key := fmt.Sprintf("banned:%s", ip)
	isBanned, err := s.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return isBanned > 0, nil
}

func (s *storage) SetBanned(ctx context.Context, ip string, banDuration time.Duration) error {
	ctx, span := s.tracer.Start(ctx, "redis::SetBanned")
	defer span.End()

	key := fmt.Sprintf("banned:%s", ip)
	return s.rdb.SetEx(ctx, key, true, banDuration).Err()
}

func (s *storage) SetCaptchaHash(ctx context.Context, id, hash string, ttl time.Duration) error {
	ctx, span := s.tracer.Start(ctx, "redis::SetCaptchaHash")
	defer span.End()

	key := fmt.Sprintf("captcha:%s", id)
	return s.rdb.SetEx(ctx, key, hash, ttl).Err()
}

func (s *storage) PopCaptchaHash(ctx context.Context, id string) (string, error) {
	ctx, span := s.tracer.Start(ctx, "redis::PopCaptchaHash")
	defer span.End()

	key := fmt.Sprintf("captcha:%s", id)
	hash, err := s.rdb.GetDel(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return hash, nil
}
