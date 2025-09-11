/*
The message reader is the file to read the message 

*/

package messagereader


import (
	"log"
	"fmt"
	"github.com/gorilla/websocket"
)


func ReadMessage(conn *websocket.Conn ) {

	// read the message
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
        	log.Println(err)
            return
        }
    
		// print out that message for clarity
        fmt.Println(string(p))
	}


}
