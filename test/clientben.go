// works in mac 

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

// Message defines the structure sent to the server
type Message struct {
	Hash      string `json:"hash"`
	Sender    string `json:"sender"`
	Recipient string `json:"receiver"`
	Message   string `json:"message"`
}

func main() {
	// hardcode for now
	chatHash := "DZ5uqxGii"
	sender := "Ben"
	recipient := "Bob" // 👈 change this to whoever you want to send messages to

	conn, _, err := websocket.DefaultDialer.Dial(
		fmt.Sprintf("ws://localhost:8080/chat-server/ws?hash=%s&user=%s", chatHash, sender),
		nil,
	)
	if err != nil {
		log.Fatal("Dial error: ", err)
	}
	defer conn.Close()

	fmt.Println("Connected to websocket Server")
	fmt.Println("Type the message and press Enter to send")

	// Goroutine to read incoming messages
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			fmt.Println("Server message:", string(msg))
		}
	}()

	// Main loop: read terminal input and send JSON
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()

		// Build the JSON message
		outMsg := Message{
			Hash:      chatHash,
			Sender:    sender,
			Recipient: recipient,
			Message:   text,
		}

		// Marshal to JSON
		data, err := json.Marshal(outMsg)
		if err != nil {
			log.Println("JSON marshal error:", err)
			continue
		}

		// Send to server
		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}
