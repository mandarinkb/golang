package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mandarinkb/go-start-bot-with-kafka/utils"
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

func KafkaConn() sarama.SyncProducer {
	client := []string{kafkaHost}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(client, config)
	if err != nil {
		fmt.Println("KafkaConn close, err:", err)
	}
	return producer
}
func KafkaAdminConn() sarama.ClusterAdmin {
	client := []string{kafkaHost}
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	clusterAdmin, err := sarama.NewClusterAdmin(client, config)
	if err != nil {
		fmt.Println("KafkaAdminConn close, err:", err)
	}
	return clusterAdmin
}

func Conn() (*sql.DB, error) {
	return sql.Open(driverName, datasourceName)
}
