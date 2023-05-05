package redsync

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/mandarinkb/go-redsync/config"
	goredislib "github.com/redis/go-redis/v9"
)

var rs *redsync.Redsync

func NewClient(cfg config.Redis) {
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     cfg.Addr,
		Password: cfg.Passwrod, // set password
		DB:       cfg.DBIndex,  // use default DB
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs = redsync.New(pool)
}

func NewMutex(key string) *redsync.Mutex {
	return rs.NewMutex(key)
}
