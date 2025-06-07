package util

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var CTX = context.Background()
var RDB = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	DB:       0,
	Password: "mandarinkb",
})
