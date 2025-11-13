/*

Get the message data for the user on login
*/

package messagefetch 

import (
	"context"
	"encoding/json"
	//"fmt"
	"log"
	"net/http"

	//"chat-service/internal/chat/session"
	"chat-service/internal/config"
	pgsqlcrud "chat-service/internal/storage/pgsql/crud"


	
)




// LoginRequest represents the expected JSON body for login requests.
type MessageFetchRequest struct {
	ChatID     string `json:"Hash"`     // Chat session hash
	UserName string `json:"UserName"` // Username trying to log in
	MessageID int `json:"MessageID"`  // added for messageid 
}


type MessagesResponse struct {
    ChatID     string    `json:"chatId"`
    Messages   []config.OutgoingMessage `json:"messages"`
    Pagination struct {
        HasMore    bool `json:"hasMore"`
        NextCursor int  `json:"nextCursor"`
    } `json:"pagination"`
}

// ErrorResponse represents a standardized error response.
type ErrorResponse struct {
	Status  string `json:"status"`  // Always "error"
	Code    int    `json:"code"`    // HTTP status code
	Message string `json:"message"` // Human-readable error message
}

// writeError writes a standardized JSON error response.
func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  "error",
		Code:    code,
		Message: msg,
	})
}


type MessageFetchResponse struct {

	message string `json:"message"`
}


// UserMessageFetch handles fetching messages for a user in a chat session.
//
// Flow:
//   1. Enforce POST method.
//   2. Decode request payload { chatId, userName }.
//   3. Validate that the user belongs to the chat (via login table).
//   4. Fetch messages for the chat from Postgres.
//   5. Return messages as JSON with pagination metadata.
func UserMessageFetch(w http.ResponseWriter, r *http.Request) {
	// Enforce POST method
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	// Decode request JSON
	var data MessageFetchRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	log.Printf("Message fetch attempt: user=%s chat=%s", data.UserName, data.ChatID)

	ctx := context.Background()
	pool := config.GlobalDbConn.PgsqlConn

	// Retrieve login record (chat participants)
	retrievedLogin, err := pgsqlcrud.GetLoginData(ctx, config.LoginTable, pool, data.ChatID)

	if err != nil {
		writeError(w, http.StatusUnauthorized, "Chat not found for given hash")
		log.Printf("Chat not found: hash=%s err=%v", data.ChatID, err)
		return
	}

	// Verify that the user belongs to this chat
	if data.UserName != retrievedLogin.UserOne && data.UserName != retrievedLogin.UserTwo {
		writeError(w, http.StatusUnauthorized, "Invalid username for this chat")
		log.Printf("Invalid username %s for chat %s", data.UserName, data.ChatID)
		return
	}

	// define the messagedata
	var messageData []config.MessageData

	// add the check for the messageid if present in the json 
	if data.MessageID == 0 {

	// Fetch messages from DB
	messageData = pgsqlcrud.GetMessageData(ctx, config.MessageTable, pool, data.ChatID, data.UserName )

	} else {
	
	// Get the message using the messageID
	messageData = pgsqlcrud.GetMessageDataID(ctx, config.MessageTable, pool, data.ChatID, data.UserName , data.MessageID)
	
	}
	
	// Convert DB model → OutgoingMessage model
	outMessages := make([]config.OutgoingMessage, len(messageData))

	for i, m := range messageData {
		outMessages[i] = config.OutgoingMessage{
			MessageID: int64(m.MessageID), 
			ChatID:    m.ChatID,
			Sender:    m.Sender,
			Receiver:  m.Receiver,
			Message:   m.Message,
			Timestamp: m.Timestamp,
		}
	}

	// Build final response
	resp := MessagesResponse{
		ChatID:   data.ChatID,
		Messages: outMessages,
		Pagination: struct {
			HasMore    bool `json:"hasMore"`
			NextCursor int  `json:"nextCursor"`
		}{
			HasMore:    false,
			NextCursor: 0,
		},
	}

	
	// Build response
	/*
	resp := MessagesResponse{
		ChatID:   data.ChatID,
		Messages: messageData, // <-- ensure pgsqlcrud.GetMessageData returns []Message
	}
	// Placeholder pagination (later implement limit/offset or cursor)
	resp.Pagination.HasMore = false
	resp.Pagination.NextCursor = 0
	*/

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("JSON encode error: %v", err)
	}

}