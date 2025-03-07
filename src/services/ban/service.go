package ban

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrTooManyRequests = errors.New("too much requests")
)

type Service struct {
	r *redis.Client
}

func NewService(r *redis.Client) *Service {
	return &Service{
		r: r,
	}
}

func (s *Service) Check(ctx context.Context, ip string) error {
	key := fmt.Sprintf("banned:%s", ip)
	isBanned, err := s.r.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if isBanned == 1 {
		return ErrTooManyRequests
	}

	key = fmt.Sprintf("reqs:%s", ip)
	pipe := s.r.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, 10*time.Second)
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	if incr.Val() > 10 {
		key := fmt.Sprintf("banned:%s", ip)
		s.r.Set(ctx, key, 1, 30*time.Minute)
		return ErrTooManyRequests
	}

	return nil
}
