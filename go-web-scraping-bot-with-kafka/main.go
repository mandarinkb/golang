package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/controller"
	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/database"
	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/utils"
)

func main() {
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	// close log
	defer logger.Sync()

	logger.Info("[web scrapping bot] start", utils.Url("-"),
		utils.User("-"), utils.Type(utils.TypeBot))
	fmt.Println(time.Now(), "web scrapping bot start")

	consumer := database.KafkaConsumerConn()
	defer consumer.Close()

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient("bot-group", consumer)
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
		topics := []string{"detail-url"}
		handler := controller.ConsumerGroupHandler{}
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}
