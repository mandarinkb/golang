package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

var (
	user = User{}
)

type exampleConsumerGroupHandler struct{}

func (exampleConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (exampleConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		fmt.Println("Partition: ", msg.Partition)
		fmt.Println("Offset: ", msg.Offset)
		fmt.Println("Key: ", string(msg.Key))

		if err := json.Unmarshal(msg.Value, &user); err != nil {
			fmt.Println(err)
		}
		fmt.Println("Value Username: ", user.Username)
		fmt.Println("Value Email: ", user.Email)
		fmt.Println()
		session.MarkMessage(msg, "")
	}
	return nil
}
func main() {
	// Init config, specify appropriate version
	config := sarama.NewConfig()
	config.Version = sarama.V3_0_0_0
	config.Consumer.Return.Errors = true

	// Start with a client
	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient("my-group", client)
	if err != nil {
		panic(err)
	}
	defer group.Close()

	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		topics := []string{"project-final"}
		handler := exampleConsumerGroupHandler{}

		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}
