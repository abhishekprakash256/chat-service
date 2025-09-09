/*
make the session


type RedisSessionData struct {
	ChatID      string
	Sender       string
	Reciever	string
	LastSeen    time.Time
	WSConnected int
	Notify      int
}

session:abc123:Abhi → {
  chat_id: abc123,
  sender: Abhi,
  reciver: Anny,
  last_seen: 2025-07-08T20:00:00,
  ws_connected: 1,
  notify: 0
}

*/

package session 

import (

	"fmt"
	//"chat-service/api/db_connector"
	"time"
	"context"

	wshandler "chat-service/api/ws"

	"chat-service/internal/config"
	rediscrud "chat-service/internal/storage/redis/crud"

)



// SaveSession stores session data in Redis.
//
// Params:
//   - hash: Chat session identifier
//   - sender: User initiating the session
//   - receiver: Opposite user
//   - lastSeen: Last active timestamp
//   - wsStatus: WebSocket connection status (1 = connected, 0 = disconnected)
//   - notification: Notification flag (1 = on, 0 = off)
//
// Returns:
//   - true if the session was saved successfully, false otherwise
func SaveSession(hash string, sender string, receiver string, lastSeen time.Time, wsStatus int, notification int) bool {
	
	client := config.GlobalDbConn.RedisConn

	ctx := context.Background()

	sessionData := config.RedisSessionData{
		Hash:        hash,
		Sender:      sender,
		Reciever:    receiver,
		LastSeen:    lastSeen,
		WSConnected: wsStatus,
		Notify:      notification,
	}

	// Session key format: session:<hash>:<sender>
	sessionID := fmt.Sprintf("session:%s:%s", hash, sender)

	err := rediscrud.StoreSessionData(ctx, client, sessionID, sessionData)
	
	if err != true {
		fmt.Println(" Failed to save session data:", err)
		return false
	}

	fmt.Println("Session data saved:", sessionID)
	return true

}




// strat the session
// The save session take the redis client 
// params -- > usename , hash 
// make the redis hash , take  timestamp and save the data 

func StartSession(hash string , sender string , receiver string ) {

	fmt.Println("Session started")

	//client := config.GlobalDbConn.RedisConn

	//ctx := context.Background()


	// these are testing values
	// time for testing 
	now := time.Now()

	var wsStatus int 

	wsStatus = 1 

	var notify int 

	notify = 1

	SaveSession(hash , sender , receiver , now , wsStatus ,  notify)

	// Session key format: session:<hash>:<sender>
	sessionID := fmt.Sprintf("session:%s:%s", hash, sender)

	// start the ws end point
	wshandler.WsHandler( sessionID )

	fmt.Println("Session Done")

}



