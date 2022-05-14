package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

var (
	user     User
	wsConn   *websocket.Conn
	username string
)

type exampleConsumerGroupHandler struct{}

func (exampleConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (exampleConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		fmt.Println("Partition: ", msg.Partition)
		fmt.Println("Offset: ", msg.Offset)
		fmt.Println("Key: ", string(msg.Key))
		fmt.Println("Value: ", string(msg.Value))
		// fmt.Println(username)

		if err := json.Unmarshal(msg.Value, &user); err != nil {
			fmt.Println(err)
		}
		if err := wsConn.WriteJSON(user); err != nil {
			return nil
		}
		fmt.Println()
		session.MarkMessage(msg, "")

	}
	return nil
}

// use default options
var upgrader = websocket.Upgrader{}

func socketHandler(w http.ResponseWriter, r *http.Request) {

	// allowed all origins
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
	var err error
	wsConn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	defer wsConn.Close()

	// Init config, specify appropriate version
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Start with a client
	client, err := sarama.NewClient([]string{"localhost:9093"}, config)
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
	username = mux.Vars(r)["username"]
	for {
		topics := []string{"project-final"}
		handler := exampleConsumerGroupHandler{}
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	route := mux.NewRouter()
	fmt.Println("websocket server port 8989 running...")
	route.HandleFunc("/socket/{username}", socketHandler)
	http.ListenAndServe(":8989", route)
}
