/*
The message reader is the file to read the message 

*/

package messagereader


import (
	"log"
	"fmt"
	"github.com/gorilla/websocket"
	"chat-service/internal/config"
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
	}


}
