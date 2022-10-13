package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"

	"github.com/mandarinkb/go-viper-config-and-response/config"
	"github.com/mandarinkb/go-viper-config-and-response/logger"
)

var (
	RedisClientData RedisClient
)

type RedisClient struct {
	redisClient *redis.Client
}

func NewClient(cfg config.Redis) {
	rd := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Passwrod, // set password
		DB:       cfg.DBIndex,  // use default DB
	})
	if err := rd.Ping(context.Background()).Err(); err != nil {
		logger.L().Errorf("cannot connect to redis: %+v", err)
	}
	RedisClientData = RedisClient{redisClient: rd}
}

func (r RedisClient) SaveDataOnRedis(ctx context.Context, key string, data interface{}, timeout time.Duration) error {
	js, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %v", err)
	}
	return r.redisClient.Set(ctx, key, js, timeout).Err()
}

func (r RedisClient) GetDataFromRedis(ctx context.Context, key string, dst interface{}) error {
	val, err := r.redisClient.Get(ctx, key).Result()
	if err != nil || err == redis.Nil {
		return err
	}
	if err := json.Unmarshal([]byte(val), &dst); err != nil {
		return fmt.Errorf("unable to unmarshal data: %v", err)
	}
	return nil
}

func (r RedisClient) RemoveDataOnRedis(ctx context.Context, key string) error {
	ssc := r.redisClient.Keys(ctx, key)
	if ssc == nil || len(ssc.Val()) == 0 {
		return redis.Nil
	}
	keys := ssc.Val()
	for _, v := range keys {
		if _, err := r.redisClient.Del(ctx, v).Result(); err != nil {
			return err
		}
	}
	return nil
}
