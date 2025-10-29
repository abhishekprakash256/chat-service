/*
make the session

//can have sessionID unique -- chnages
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
//   - chatID: Chat session identifier
//   - sender: User initiating the session
//   - receiver: Opposite user
//   - lastSeen: Last active timestamp
//   - wsStatus: WebSocket connection status (1 = connected, 0 = disconnected)
//   - notification: Notification flag (1 = on, 0 = off)
//
// Returns:
//   - error if the operation fails, nil otherwise

// add sessionID here now -- chnages 
func SaveSession(chatID string, sessionID string , sender string, receiver string, lastSeen time.Time, wsStatus int, notification int) error {
	client := config.GlobalDbConn.RedisConn
	ctx := context.Background()

	sessionData := config.RedisSessionData{
		ChatID:      chatID,
		SessionID:   sessionID,
		Sender:      sender,
		Reciever:    receiver,
		LastSeen:    lastSeen,
		WSConnected: wsStatus,
		Notify:      notification,
	}

	// Session key format: session:<chatID>:<sender>
	sessionKey := fmt.Sprintf("session:%s:%s:%s", chatID, sender, sessionID)

	ok := rediscrud.StoreSessionData(ctx, client, sessionKey , sessionData)
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
//   - chatID: Chat session identifier
//   - sender: Username of the client who initiated login
//
// Flow:
//   1. Retrieves receiver from PostgreSQL login table.
//   2. Creates a Redis session entry.
//   3. Starts a heartbeat goroutine to monitor connection health.


//StartSession(conn , chatID , sessionID , sender )
func StartSession(conn *websocket.Conn, chatID string , sessionID string , sender string) {
	
	fmt.Println("Session started")

	ctx := context.Background()
	pool := config.GlobalDbConn.PgsqlConn

	// Retrieve login record from DB
	retrievedLogin, err := pgsqlcrud.GetLoginData(ctx, config.LoginTable, pool, chatID)
	if err != nil {
		log.Printf("Failed to get login data for chatID %s: %v", chatID, err)
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
	if err := SaveSession(chatID, sessionID , sender, receiver, now, 1, 1); err != nil {
		log.Printf("Could not save session for %s: %v", sender, err)
		return
	}

	// Start heartbeat
	go startHeartbeat(conn, sender, receiver, chatID ,sessionID )

	//can start the message hub here as well 
}


// startHeartbeat runs a periodic ping to check if the client is alive.
// If the heartbeat fails, the session is marked as disconnected in Redis.
//
// Params:
//   - conn: WebSocket connection
//   - sender: Username of the session owner
//   - receiver: Opposite user in the chat
//   - chatID: Chat session identifier

func startHeartbeat(conn *websocket.Conn, sender, receiver, chatID, sessionID string) {
    ticker := time.NewTicker(30 * time.Second) // send ping every 30s
    defer ticker.Stop()

    wsKey := fmt.Sprintf("session:%s:%s", chatID, sender)
    sessionKey := fmt.Sprintf("session:%s:%s:%s", chatID, sender, sessionID)

    for {
        <-ticker.C

        // Send ping frame
        if err := conn.WriteControl(
            websocket.PingMessage,
            []byte("ping"),
            time.Now().Add(5*time.Second),
        ); err != nil {

            log.Printf("Heartbeat failed for %s and %s: %v", sessionKey, wsKey , err)
            
            /* Testing
            // --- Cleanup begins ---
            conn.Close()

            // remove from ClientsWsMapper safely
            config.ClientsWsMapper.Lock()
            if sessions, ok := config.ClientsWsMapper.Data[wsKey]; ok {
                delete(sessions, sessionKey)
                if len(sessions) == 0 {
                    delete(config.ClientsWsMapper.Data, wsKey)
                }
            }
            config.ClientsWsMapper.Unlock()

            // mark session as disconnected in Redis
            now := time.Now()
            if err := SaveSession(chatID, sessionID , sender, receiver, now, 0, 0); err != nil {
                log.Printf("Failed to update session status for %s: %v", sessionKey, err)
            }
            

            log.Printf("Session cleaned up for %s", sessionKey)
            return

            */
        }


        log.Printf("Heartbeat OK for %s", sessionKey)
    }
}



// StopSession cleanly closes a user's session.
//
// Params:
//   - conn: active WebSocket connection
//   - chatID: chat session identifier
//   - sender: the user logging out / disconnecting
//   - receiver: the opposite user
func StopSession(conn *websocket.Conn, chatID, sessionID, sender, receiver string) {
    // Create Redis and map keys
    wsKey := fmt.Sprintf("session:%s:%s", chatID, sender)
    sessionKey := fmt.Sprintf("session:%s:%s:%s", chatID, sender, sessionID)

    // --- Step 1: Gracefully close the connection ---
    if conn != nil {
        _ = conn.WriteMessage(
            websocket.CloseMessage,
            websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Session closed"),
        )
        _ = conn.Close()
    }

    // --- Step 2: Thread-safe cleanup from ClientsWsMapper ---
    config.ClientsWsMapper.Lock()
    if sessionMap, exists := config.ClientsWsMapper.Data[wsKey]; exists {
        delete(sessionMap, sessionKey)
        if len(sessionMap) == 0 {
            delete(config.ClientsWsMapper.Data, wsKey)
        }
    }
    config.ClientsWsMapper.Unlock()

    log.Printf("Stopped and removed session: %s", sessionKey)

    // --- Step 3: Mark session inactive in Redis ---
    now := time.Now()
    if err := SaveSession(chatID, sessionID , sender, receiver, now, 0, 0); err != nil {
        log.Printf("Failed to update session status for %s: %v", sessionKey, err)
    }
}


func RemoveClient(chatID, user, sessionID string) {
    wskey := fmt.Sprintf("session:%s:%s", chatID, user)
    sessionKey := fmt.Sprintf("session:%s:%s:%s", chatID, user, sessionID)

    config.ClientsWsMapper.Lock()
    defer config.ClientsWsMapper.Unlock()

    if sessions, ok := config.ClientsWsMapper.Data[wskey]; ok {
        delete(sessions, sessionKey)
        if len(sessions) == 0 {
            delete(config.ClientsWsMapper.Data, wskey)
        }
    }
}
