package guard

import (
	"context"
	"time"
)

type Storage interface {
	IncIpRate(ctx context.Context, ip string, window time.Duration) (int, error)
	IsIpBanned(ctx context.Context, ip string) (bool, error)
	SetIpBanned(ctx context.Context, ip string, banDuration time.Duration) error
	SetCaptchaHash(ctx context.Context, id, hash string, ttl time.Duration) error
	PopCaptchaHash(ctx context.Context, id string) (string, error)
}
