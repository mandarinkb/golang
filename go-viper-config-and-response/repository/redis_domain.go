package repository

import (
	"context"
	"time"
)

type RedisRepository interface {
	SaveRedis(ctx context.Context, key string, data interface{}, timeout time.Duration) error
	GetRedis(ctx context.Context, key string, dst interface{}) error
	RemoveRedis(ctx context.Context, key string) error
}
