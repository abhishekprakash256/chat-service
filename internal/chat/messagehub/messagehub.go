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
	//"chat-service/internal/chat/messagestore"


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
            incoming.Sender, incoming.Receiver, incoming.ChatID)

        // Construct session key for the intended recipient
        wsKey := fmt.Sprintf("session:%s:%s", incoming.ChatID, incoming.Receiver)

        config.ClientsWsMapper.RLock()
        sessions, ok := config.ClientsWsMapper.Data[wsKey]
        config.ClientsWsMapper.RUnlock()

        // Look up the recipient connections
        if !ok || len(sessions) == 0 {
            log.Printf("Recipient %s not connected for chat %s", incoming.Receiver, incoming.ChatID)
            continue
        }

       // Deliver message to all active sessions for the recipient
        config.ClientsWsMapper.Lock()
        for sessionID, conn := range sessions {
            if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {

                log.Printf("Failed to deliver to %s (session %s): %v", wsKey, sessionID, err)
                conn.Close()
                delete(sessions, sessionID)
                continue
            }
        }

        // Keep only alive connections
        if len(sessions) == 0 {
            delete(config.ClientsWsMapper.Data, wsKey)
        } else {
            config.ClientsWsMapper.Data[wsKey] = sessions
        }
        config.ClientsWsMapper.Unlock()
        

    }
}