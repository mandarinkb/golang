package main

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func main() {
	//addresses of available kafka client
	client := []string{"localhost:9092"}
	//setup relevant config info
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(client, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	defer producer.Close()

	user := User{
		Username: "man",
		Email:    "man@gmail.com",
	}

	json, _ := json.Marshal(user)

	topic := "project-final"
	msg := string(json)
	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("start-bot"),
		Value: sarama.StringEncoder(msg),
	}

	pid, offset, err := producer.SendMessage(message)
	if err != nil {
		fmt.Println("send message failed,", err)
		return
	}

	fmt.Printf("partition:%v offset:%v\n", pid, offset)
}
