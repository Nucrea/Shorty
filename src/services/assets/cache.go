package assets

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"go.opentelemetry.io/otel/trace"
)

type cache struct {
	rdb    *redis.Client
	tracer trace.Tracer
}

func (c *cache) PutMetadata(ctx context.Context, meta AssetMetadataDTO) error {
	_, span := c.tracer.Start(ctx, "redis::PutMetadata")
	defer span.End()

	bytes, err := msgpack.Marshal(meta)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("asset:%s", meta.Id)
	_, err = c.rdb.SetEx(ctx, key, bytes, time.Hour).Result()
	return err
}

func (c *cache) GetMetadata(ctx context.Context, id string) (*AssetMetadataDTO, error) {
	_, span := c.tracer.Start(ctx, "redis::GetMetadata")
	defer span.End()

	key := fmt.Sprintf("asset:%s", id)
	bytes, err := c.rdb.GetEx(ctx, key, time.Hour).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	meta := &AssetMetadataDTO{}
	if err := msgpack.Unmarshal(bytes, &meta); err != nil {
		return nil, err
	}

	return meta, nil
}
