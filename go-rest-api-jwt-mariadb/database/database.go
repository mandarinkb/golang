package database

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
)

type database struct{}

func NewDatabase() database {
	return database{}
}

func (database) RedisConn() *redis.Client {
	return redis.NewClient(&redis.Options{
		// Addr:     "localhost:6379",
		Addr:     "redis:6379",
		Password: "mandarinkb", // set password
		DB:       0,            // use default DB
	})
}

func (database) Conn() (*sql.DB, error) {
	// return sql.Open("mysql", "root:mandarinkb@tcp(127.0.0.1)/TEST_DB?charset=utf8")
	return sql.Open("mysql", "root:mandarinkb@tcp(mariadb)/TEST_DB?charset=utf8")
}
