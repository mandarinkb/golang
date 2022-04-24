package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mandarinkb/go-fetch-url-bot-with-kafka/utils"
)

var (
	kafkaHost      string
	driverName     string
	datasourceName string
)

func init() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	kafkaHost = config.KafkaHost
	driverName = config.DriverName
	datasourceName = config.DatasourceName
}

func KafkaProducerConn() sarama.SyncProducer {
	client := []string{kafkaHost}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(client, config)
	if err != nil {
		fmt.Println("KafkaConn close, err:", err)
	}
	return producer
}

func KafkaConsumerConn() sarama.Client {
	client := []string{kafkaHost}
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewClient(client, config)
	if err != nil {
		panic(err)
	}
	return consumer
}
func Conn() (*sql.DB, error) {
	return sql.Open(driverName, datasourceName)
}
