/*
The file to router the message to the user for which messasge is send 
*/

package messagehub


import (

	"chat-service/internal/config"
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"fmt"
	


)




// HandleMessages listens on the broadcast channel and routes messages
// to the correct WebSocket clients based on chat hash and recipient username.
// HandleMessages listens on the broadcast channel and routes messages
// to the correct WebSocket clients based on chat hash and recipient username.
func HandleMessages() {
    
    for {
        // Grab next message from broadcast channel
        msg := <-config.BroadCast

        // Unmarshal the incoming JSON into a struct
        var incoming config.IncomingMessage

        if err := json.Unmarshal(msg, &incoming); err != nil {
            log.Println("Invalid message format:", err)
            continue
        }

        log.Printf("Routing message from %s to %s in chat %s",
            incoming.Sender, incoming.Receiver, incoming.Hash)

        // Construct session key for the intended recipient
        chatID := fmt.Sprintf("session:%s:%s", incoming.Hash, incoming.Receiver)

        // Look up the recipient connections
        conns, ok := config.ClientsWsMapper[chatID]
        if !ok || len(conns) == 0 {
            log.Printf("Recipient %s not connected for chat %s",
                incoming.Receiver, incoming.Hash)
            continue
        }

        // Deliver message to all active connections for that recipient
        aliveConns := []*websocket.Conn{}
        for _, conn := range conns {
            if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
                log.Printf("Failed to deliver to %s: %v", chatID, err)
                conn.Close()
                continue // skip dead connection
            }
            aliveConns = append(aliveConns, conn)
        }

        // Keep only alive connections
        if len(aliveConns) > 0 {
            config.ClientsWsMapper[chatID] = aliveConns
        } else {
            delete(config.ClientsWsMapper, chatID)
        }
    }
}