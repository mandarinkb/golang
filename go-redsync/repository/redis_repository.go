package repository

import (
	"context"
	"time"

	"github.com/mandarinkb/go-redsync/database"
	"github.com/mandarinkb/go-redsync/logger"
)

type redisRepository struct {
	log *logger.Logger
}

func NewRedisRepository() RedisRepository {
	return &redisRepository{
		log: logger.L().Named("redis repository"),
	}
}

func (r redisRepository) SaveRedisNoneExpiry(ctx context.Context, key string, src interface{}) error {
	return database.RedisClientData.SaveDataOnRedisNoneExpiry(ctx, key, src)
}

func (r redisRepository) SaveRedisWithTTL(ctx context.Context, key string, src interface{}, timeout time.Duration) error {
	return database.RedisClientData.SaveDataOnRedisWithTTL(ctx, key, src, timeout)
}

func (r redisRepository) GetRedis(ctx context.Context, key string, dst interface{}) error {
	return database.RedisClientData.GetDataFromRedis(ctx, key, dst)
}

func (r redisRepository) RemoveRedis(ctx context.Context, key string) error {
	return database.RedisClientData.RemoveDataOnRedis(ctx, key)
}

func (r redisRepository) PushRedis(ctx context.Context, key string, src interface{}) error {
	return database.RedisClientData.PushDataOnRedis(ctx, key, src)
}

func (r redisRepository) PopRedis(ctx context.Context, key string, dst interface{}) error {
	return database.RedisClientData.PopDataFromRedis(ctx, key, dst)
}

func (r redisRepository) ProducerRedis(ctx context.Context, key string, src interface{}) error {
	return database.RedisClientData.PublishDataOnRedis(ctx, key, src)
}

func (r redisRepository) ConsumerRedis(ctx context.Context, key string, dst interface{}) error {
	return database.RedisClientData.SubscribeDataFromRedis(ctx, key, dst)
}
