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
	"log"
	"time"
	"context"
	"github.com/gorilla/websocket"

	//wshandler "chat-service/api/ws"

	"chat-service/internal/config"
	rediscrud "chat-service/internal/storage/redis/crud"
	pgsqlcrud "chat-service/internal/storage/pgsql/crud"

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

func StartSession(conn *websocket.Conn ,  hash string , sender string ) {

	fmt.Println("Session started")

	ctx := context.Background()
	pool := config.GlobalDbConn.PgsqlConn

	// Retrieve login record from DB
	retrievedLogin, _ := pgsqlcrud.GetLoginData(ctx, "login", pool, hash)

	var receiver string

	// Check if username matches registered users
	if sender == retrievedLogin.UserOne {

		receiver = retrievedLogin.UserTwo

	} else  {

		receiver = retrievedLogin.UserOne
	
	} 

	//sessionID := fmt.Sprintf("session:%s:%s", hash, sender)


	go startHeartbeat(conn, sender , receiver , hash  )

}


func startHeartbeat(conn *websocket.Conn, sender string , receiver string , hash string) {

    ticker := time.NewTicker(30 * time.Second) // every 30s
    defer ticker.Stop()

	sessionID := fmt.Sprintf("session:%s:%s", hash, sender)

    for {
        <-ticker.C

        // send ping
        if err := conn.WriteControl(
            websocket.PingMessage,
            []byte("ping"),
            time.Now().Add(5*time.Second),
        ); err != nil {
            log.Printf("Heartbeat failed for %s: %v", sessionID, err)


            // remove from ClientsWsMapper
            delete(config.ClientsWsMapper, sessionID)

			now := time.Now()

            // update Redis: ws_status = 0
            SaveSession(hash , sender , receiver , now , 1,  0)

            return
        }

        log.Printf("Heartbeat OK for %s", sessionID)
    }
}


