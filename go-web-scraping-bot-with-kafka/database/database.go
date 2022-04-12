package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/utils"
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
	//addresses of available kafka client
	client := []string{kafkaHost}
	//setup relevant config info
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(client, config)
	if err != nil {
		fmt.Println("KafkaConn close, err:", err)
	}
	return producer
}
func KafkaProducerAdminConn() sarama.ClusterAdmin {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	clusterAdmin, err := sarama.NewClusterAdmin([]string{kafkaHost}, config)
	if err != nil {
		fmt.Println("KafkaAdminConn close, err:", err)
	}
	return clusterAdmin
}

func KafkaConsumerConn() sarama.Client {
	config := sarama.NewConfig()
	config.Version = sarama.V3_0_0_0
	config.Consumer.Return.Errors = true
	client, err := sarama.NewClient([]string{"localhost:9093"}, config)
	if err != nil {
		panic(err)
	}
	return client
}
func Conn() (*sql.DB, error) {
	return sql.Open(driverName, datasourceName)
}
