package database

import (
	"database/sql"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mandarinkb/go-api-project-final/utils"
)

type database struct{}

func NewDatabase() database {
	return database{}
}

var redisHost string
var redisPassword string
var driverName string
var datasourceName string

func init() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	redisHost = config.RedisHost
	redisPassword = config.RedisPassword
	driverName = config.DriverName
	datasourceName = config.DatasourceName
}

func (database) RedisConn() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword, // set password
		DB:       0,             // use default DB
	})
}

func (database) Conn() (*sql.DB, error) {
	return sql.Open(driverName, datasourceName)
}
