package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// use default options
var upgrader = websocket.Upgrader{}

type data struct {
	Item string `json:"item"`
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func socketHandler(w http.ResponseWriter, r *http.Request) {
	// allowed all origins
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	defer conn.Close()
	// ##### ส่ง response ทุกๆ 2 วินาที #####
	for {
		// random string 10 letter
		obj := data{
			Item: randomStr(10),
		}

		if err = conn.WriteJSON(obj); err != nil {
			return
		}
		time.Sleep(2 * time.Second)
	}
	// ##### รับข้อมูลจาก request แล้วส่งกลับ response #####
	// for {
	// 	// Read message from browser
	// 	msgType, msg, err := conn.ReadMessage()
	// 	if err != nil {
	// 		return
	// 	}

	// 	// Print the message to the console
	// 	fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
	// 	// Write message back to browser
	// 	if err = conn.WriteMessage(msgType, msg); err != nil {
	// 		return
	// 	}
	// }

}

func main() {
	fmt.Println("websocket server port 8989 running...")
	http.HandleFunc("/socket", socketHandler)
	http.ListenAndServe(":8989", nil)
}
