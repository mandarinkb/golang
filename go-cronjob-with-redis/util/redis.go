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

	// Redis config
	addr := GetEnv("REDIS_ADDR", "localhost:6379")
	pass := GetEnv("REDIS_PASSWORD", "mandarinkb")
	db := 0

	RDB = redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    pass,
		DB:          db,
		ReadTimeout: 2 * time.Second,
	})
	if err := RDB.Ping(CTX).Err(); err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}

	Logger.Info("✅ Redis connected", slog.String("addr", addr))
}

func GetEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
