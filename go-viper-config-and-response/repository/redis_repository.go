package repository

import (
	"context"
	"time"

	"github.com/mandarinkb/go-viper-config-and-response/database"
	"github.com/mandarinkb/go-viper-config-and-response/logger"
)

type redisRepository struct {
	log *logger.Logger
}

func NewRedisRepository() RedisRepository {
	return &redisRepository{
		log: logger.L().Named("redis repository"),
	}
}

func (r redisRepository) SaveRedis(ctx context.Context, key string, data interface{}, timeout time.Duration) error {
	return database.RedisClientData.SaveDataOnRedis(ctx, key, data, timeout)
}

func (r redisRepository) GetRedis(ctx context.Context, key string, dst interface{}) error {
	return database.RedisClientData.GetDataFromRedis(ctx, key, dst)
}

func (r redisRepository) RemoveRedis(ctx context.Context, key string) error {
	return database.RedisClientData.RemoveDataOnRedis(ctx, key)
}
