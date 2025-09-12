// works in mac 
/*
making the go client
*/

package main 

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"github.com/gorilla/websocket"

)

func main() {

	conn , _ , err := websocket.DefaultDialer.Dial("ws://localhost:8080/chat-server/ws?hash=XWVU7wbbr&user=Bob", nil)
	//conn , _ , err := websocket.DefaultDialer.Dial("ws://localhost:8080/chat-server/ws?hash=XWVU7wbbr&user=Ben", nil)


	if err != nil {

		log.Fatal("Dial error: ", err)
	}
	
	defer conn.Close()

	fmt.Println("Connected to websocket Server")

	fmt.Println("Type the message and press Enter to send ")

	// start a goroutine to continouly read the message from the server 

	go func() {

		for {

			_, msg , err := conn.ReadMessage()

			if err != nil {
				log.Println("Read error: " , err )
			}

			fmt.Println("Server message :", string(msg))
		}
	}()

	// Main loop : read the user inpur from the terminal and send to the server

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		text := scanner.Text()

		err := conn.WriteMessage(websocket.TextMessage , []byte(text) )

		if err != nil {
			log.Println("Write error: " , err )

			return 
		}
	}


}