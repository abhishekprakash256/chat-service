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
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"chat-service/internal/config"
	rediscrud "chat-service/internal/storage/redis/crud"
	pgsqlcrud "chat-service/internal/storage/pgsql/crud"
)


// SaveSession stores chat session metadata in Redis.
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
//   - error if the operation fails, nil otherwise
func SaveSession(hash string, sender string, receiver string, lastSeen time.Time, wsStatus int, notification int) error {
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

	ok := rediscrud.StoreSessionData(ctx, client, sessionID, sessionData)
	if !ok {
		return fmt.Errorf("failed to save session data for %s", sessionID)
	}

	log.Println("Session data saved:", sessionID)
	return nil
}


// StartSession initializes a new chat session and starts the heartbeat.
//
// Params:
//   - conn: Active WebSocket connection
//   - hash: Chat session identifier
//   - sender: Username of the client who initiated login
//
// Flow:
//   1. Retrieves receiver from PostgreSQL login table.
//   2. Creates a Redis session entry.
//   3. Starts a heartbeat goroutine to monitor connection health.
func StartSession(conn *websocket.Conn, hash string, sender string) {
	fmt.Println("Session started")

	ctx := context.Background()
	pool := config.GlobalDbConn.PgsqlConn

	// Retrieve login record from DB
	retrievedLogin, err := pgsqlcrud.GetLoginData(ctx, "login", pool, hash)
	if err != nil {
		log.Printf("Failed to get login data for hash %s: %v", hash, err)
		return
	}

	// Determine receiver
	var receiver string
	if sender == retrievedLogin.UserOne {
		receiver = retrievedLogin.UserTwo
	} else {
		receiver = retrievedLogin.UserOne
	}

	// Save initial session
	now := time.Now()
	if err := SaveSession(hash, sender, receiver, now, 1, 1); err != nil {
		log.Printf("Could not save session for %s: %v", sender, err)
		return
	}

	// Start heartbeat
	go startHeartbeat(conn, sender, receiver, hash)

	//can start the message hub here as well 
}


// startHeartbeat runs a periodic ping to check if the client is alive.
// If the heartbeat fails, the session is marked as disconnected in Redis.
//
// Params:
//   - conn: WebSocket connection
//   - sender: Username of the session owner
//   - receiver: Opposite user in the chat
//   - hash: Chat session identifier
func startHeartbeat(conn *websocket.Conn, sender string, receiver string, hash string) {
	ticker := time.NewTicker(30 * time.Second) // send ping every 30s
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

			// cleanup: close socket and remove from mapper
			conn.Close()

			delete(config.ClientsWsMapper, sessionID)

			// mark session as disconnected in Redis
			now := time.Now()

			if err := SaveSession(hash, sender, receiver, now, 0, 0); err != nil {

				log.Printf("Failed to update session status for %s: %v", sessionID, err)
			}
			log.Printf("Session cleaned up for %s", sessionID)

			return
		}

		log.Printf("Heartbeat OK for %s", sessionID)
	}
}
