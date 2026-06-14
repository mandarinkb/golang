package util

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var CTX = context.Background()
var RDB *redis.Client
var Logger *slog.Logger

func Init() {
	// Logger
	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(Logger)

	addr := GetEnv("REDIS_ADDR", "localhost:6379")
	pass := GetEnv("REDIS_PASSWORD", "mandarinkb")

	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,

		// Timeouts
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		// Connection pool — ป้องกัน connection exhaustion เมื่อ concurrent workflows สูง
		PoolSize:        20,
		MinIdleConns:    5,
		ConnMaxIdleTime: 5 * time.Minute,

		// Auto-retry เมื่อ connection drop ชั่วคราว
		MaxRetries:      3,
		MinRetryBackoff: 100 * time.Millisecond,
		MaxRetryBackoff: 1 * time.Second,
	})

	if err := RDB.Ping(CTX).Err(); err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}

	Logger.Info("✅ Redis connected",
		slog.String("addr", addr),
		slog.Int("pool_size", 20),
		slog.Int("max_retries", 3),
	)
}

func GetEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
