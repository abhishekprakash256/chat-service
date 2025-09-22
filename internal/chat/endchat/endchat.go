/*
The package to end the chat 
*/




package endchat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"chat-service/internal/config"
	"chat-service/internal/chat/session"
	pgsqlcrud "chat-service/internal/storage/pgsql/crud"
	rediscrud "chat-service/internal/storage/redis/crud"
)

// EndChatRequest represents the expected logout request payload.
type EndChatRequest struct {
	ChatID   string `json:"hash"`
	UserName string `json:"username"`
}

// SuccessResponse defines a standard logout success payload.
type SuccessResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse defines a standard error payload.
type ErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// writeError sends a consistent error response to the client.
func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Status:  "error",
		Code:    code,
		Message: msg,
	})
}

// writeSuccess sends a standard success response.
func writeSuccess(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(SuccessResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: msg,
	})
}

// UserEndChat handles cleanup when a user ends a chat session.
// Steps performed:
//   1. Validate request method (must be POST).
//   2. Decode request JSON and validate chat/user.
//   3. Retrieve chat participants from DB.
//   4. Close any active WebSocket connections for sender and receiver.
//   5. Delete session data from Redis for both users.
//   6. Remove related message history from PostgreSQL.
//
// On success → returns { "status": "success", "message": "Chat ended" }
func UserEndChat(w http.ResponseWriter, r *http.Request) {
	// Enforce POST method
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	// Decode request JSON
	var data EndChatRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	log.Printf("Logout attempt: %s for chat %s", data.UserName, data.ChatID)

	ctx := context.Background()
	pool := config.GlobalDbConn.PgsqlConn
	rdb := config.GlobalDbConn.RedisConn

	// Retrieve login record (chat participants)
	retrievedLogin, err := pgsqlcrud.GetLoginData(ctx, "login", pool, data.ChatID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Chat not found for given hash")
		log.Printf("Chat end failed: hash %s not found (%v)", data.ChatID, err)
		return
	}

	// Identify sender/receiver
	var sender, receiver string
	switch data.UserName {
	case retrievedLogin.UserOne:
		sender, receiver = retrievedLogin.UserOne, retrievedLogin.UserTwo
	case retrievedLogin.UserTwo:
		sender, receiver = retrievedLogin.UserTwo, retrievedLogin.UserOne
	default:
		writeError(w, http.StatusUnauthorized, "Invalid username for this chat")
		log.Printf("Invalid username %s for chat %s", data.UserName, data.ChatID)
		return
	}

	// Compose session keys
	sessionIDSender := fmt.Sprintf("session:%s:%s", data.ChatID, sender)
	sessionIDReceiver := fmt.Sprintf("session:%s:%s", data.ChatID, receiver)

	// Close all active WebSocket connections
	closeWebSockets(sessionIDSender)
	closeWebSockets(sessionIDReceiver)

	// Delete session data from Redis
	if !rediscrud.DeleteSessionData(ctx, rdb, sessionIDSender) {
		log.Printf("Warning: failed to delete redis session %s", sessionIDSender)
	}
	if !rediscrud.DeleteSessionData(ctx, rdb, sessionIDReceiver) {
		log.Printf("Warning: failed to delete redis session %s", sessionIDReceiver)
	}

	// Remove chat messages
	if err := pgsqlcrud.DeleteMessageData(ctx, config.MessageTable, pool, data.ChatID); err != nil {
		log.Printf("Warning: failed to delete message history for chat %s: %v", data.ChatID, err)
	}

	// TODO: optionally stop heartbeat goroutines for sender/receiver
	// session.StopSession(sender, data.ChatID)

	writeSuccess(w, "Chat ended successfully")
}


// closeWebSockets safely closes and removes all connections for a session ID.
func closeWebSockets(sessionID string) {
	if conns, ok := config.ClientsWsMapper[sessionID]; ok {
		for _, conn := range conns {
			// Send close frame before closing
			_ = conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Chat ended"))
			_ = conn.Close()
		}
		delete(config.ClientsWsMapper, sessionID)
		log.Printf("Closed all WebSockets for session %s", sessionID)
	}
}
