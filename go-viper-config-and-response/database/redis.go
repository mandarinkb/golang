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

func (r RedisClient) SaveDataOnRedisNoneExpiry(ctx context.Context, key string, data interface{}) error {
	return r.saveDataOnRedis(ctx, key, data, 0)
}

func (r RedisClient) SaveDataOnRedisWithTTL(ctx context.Context, key string, data interface{}, timeout time.Duration) error {
	return r.saveDataOnRedis(ctx, key, data, timeout)
}

func (r RedisClient) saveDataOnRedis(ctx context.Context, key string, data interface{}, timeout time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %v", err)
	}
	return r.redisClient.Set(ctx, key, b, timeout).Err()
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
		if err := r.redisClient.Del(ctx, v).Err(); err != nil {
			return err
		}
	}
	return nil
}
func (r RedisClient) PushDataOnRedis(ctx context.Context, key string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %v", err)
	}
	return r.redisClient.RPush(ctx, key, b).Err()
}

func (r RedisClient) PopDataFromRedis(ctx context.Context, key string, dst interface{}) error {
	val, err := r.redisClient.RPop(ctx, key).Result()
	if err != nil || err == redis.Nil {
		return err
	}
	if err := json.Unmarshal([]byte(val), &dst); err != nil {
		return fmt.Errorf("unable to unmarshal data: %v", err)
	}
	return nil
}

func (r RedisClient) PublishDataOnRedis(ctx context.Context, key string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %v", err)
	}
	return r.redisClient.RPush(ctx, key, b).Err()
}

func (r RedisClient) SubscribeDataFromRedis(ctx context.Context, key string, dst interface{}) error {
	msg, err := r.redisClient.BLPop(ctx, config.C().RedisReadTimeout, key).Result()
	if err != nil || err == redis.Nil {
		return err
	}
	for _, v := range msg {
		if v != key {
			if err := json.Unmarshal([]byte(v), &dst); err != nil {
				return fmt.Errorf("unable to unmarshal data: %v", err)
			}
		}

	}
	return nil
}
