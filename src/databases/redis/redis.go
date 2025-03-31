package redis

import (
	"context"
	"fmt"
	"shorty/src/common/metrics"
	"shorty/src/services/assets"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"go.opentelemetry.io/otel/trace"
)

func New(ctx context.Context, redisUrl string, tracer trace.Tracer, meter metrics.Meter) (*redisDb, error) {
	redisOpts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, fmt.Errorf("error parsing redis url: %w", err)
	}
	rdb := redis.NewClient(redisOpts)

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("error pinging redis: %w", err)
	}

	return &redisDb{
		rdb:    rdb,
		tracer: tracer,
		latencyHist: meter.NewHistogram(
			"redis_query_latency",
			"Latency of redis queries",
			[]float64{1, 10, 50, 100, 200},
		),
		// errorsCounter: meter.NewCounter("redis_query_errors", "Count of Redis query errors"),
	}, nil
}

type redisDb struct {
	rdb         *redis.Client
	tracer      trace.Tracer
	latencyHist metrics.Histogram
	// errorsCounter metrics.Counter
}

func (r *redisDb) observe(ctx context.Context, funcName string) func() {
	_, span := r.tracer.Start(ctx, fmt.Sprintf("redis::%s", funcName))
	defer span.End()

	start := time.Now()
	return func() {
		span.End()
		duration := time.Now().Sub(start).Milliseconds()
		r.latencyHist.ObserveWithLabel(float64(duration), funcName)
	}
}

func (r *redisDb) IncIpRate(ctx context.Context, ip string, window time.Duration) (int, error) {
	defer r.observe(ctx, "IncIpRate")()

	key := fmt.Sprintf("rate:%s", ip)

	pipe := r.rdb.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.ExpireNX(ctx, key, window)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, err
	}

	return int(incr.Val()), nil
}

func (r *redisDb) IsIpBanned(ctx context.Context, ip string) (bool, error) {
	defer r.observe(ctx, "IsIpBanned")()

	key := fmt.Sprintf("banned:%s", ip)
	isBanned, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return isBanned > 0, nil
}

func (r *redisDb) SetIpBanned(ctx context.Context, ip string, banDuration time.Duration) error {
	defer r.observe(ctx, "SetIpBanned")()

	key := fmt.Sprintf("banned:%s", ip)
	return r.rdb.SetEx(ctx, key, true, banDuration).Err()
}

func (r *redisDb) SetCaptchaHash(ctx context.Context, id, hash string, ttl time.Duration) error {
	defer r.observe(ctx, "SetCaptchaHash")()

	key := fmt.Sprintf("captcha:%s", id)
	return r.rdb.SetEx(ctx, key, hash, ttl).Err()
}

func (r *redisDb) PopCaptchaHash(ctx context.Context, id string) (string, error) {
	defer r.observe(ctx, "PopCaptchaHash")()

	key := fmt.Sprintf("captcha:%s", id)
	hash, err := r.rdb.GetDel(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (r *redisDb) PutAssetMetadata(ctx context.Context, meta assets.AssetMetadataDTO) error {
	defer r.observe(ctx, "PutAssetMetadata")()

	bytes, err := msgpack.Marshal(meta)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("asset:%s", meta.Id)
	_, err = r.rdb.SetEx(ctx, key, bytes, time.Hour).Result()
	return err
}

func (r *redisDb) GetAssetMetadata(ctx context.Context, id string) (*assets.AssetMetadataDTO, error) {
	defer r.observe(ctx, "GetAssetMetadata")()

	key := fmt.Sprintf("asset:%s", id)
	bytes, err := r.rdb.GetEx(ctx, key, time.Hour).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	meta := &assets.AssetMetadataDTO{}
	if err := msgpack.Unmarshal(bytes, &meta); err != nil {
		return nil, err
	}

	return meta, nil
}
