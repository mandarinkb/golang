package repository

import (
	"context"
	"time"
)

type RedisRepository interface {
	SaveRedisNoneExpiry(ctx context.Context, key string, src interface{}) error
	SaveRedisWithTTL(ctx context.Context, key string, src interface{}, timeout time.Duration) error
	GetRedis(ctx context.Context, key string, dst interface{}) error
	RemoveRedis(ctx context.Context, key string) error
	PushRedis(ctx context.Context, key string, src interface{}) error
	PopRedis(ctx context.Context, key string, dst interface{}) error
	ProducerRedis(ctx context.Context, key string, src interface{}) error
	ConsumerRedis(ctx context.Context, key string, dst interface{}) error
}
