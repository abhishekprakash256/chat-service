/*
The message reader is the file to read the message 

*/

package messagereader


import (
	"log"
	"fmt"
	"encoding/json"
	"github.com/gorilla/websocket"
	"chat-service/internal/config"
	"chat-service/internal/chat/messagestore"
)


func ReadMessage(conn *websocket.Conn ) {

	// read the message
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
        	log.Println(err)
            return
        }
		
		config.BroadCast <- msg
		// print out that message for clarity
        fmt.Println(string(msg))

		var incoming config.IncomingMessage

		if err := json.Unmarshal(msg, &incoming); err != nil {
			log.Println("Invalid message format:", err)
			return 
		}

		log.Printf("Routing message from %s to %s in chat %s", incoming.Sender, incoming.Receiver, incoming.Hash)

		//save the message in pgsql
		if err := messagestore.SaveMessage(incoming.Sender , incoming.Receiver , incoming.Message , incoming.Hash); err !=nil {

			log.Println("Message not saved")
			return
		}

		log.Println("Message saved succesfully in message reader")

	}


}
