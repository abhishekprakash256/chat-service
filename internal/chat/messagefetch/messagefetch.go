/*

Get the message data for the user on login
*/

package messagefetch 

import (
	"context"
	"encoding/json"
	"fmt"
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


// check for the post request 
// check for the login 
// get the message 



func UserMessageFetch( w http.ResponseWriter, r *http.Request  ) {

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
	log.Printf("Logout attempt: %s for chat %s", data.UserName, data.ChatID)

	ctx := context.Background()

	pool := config.GlobalDbConn.PgsqlConn

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

	fmt.Printf("Messages: %s\n", sessionIDSender)
	fmt.Printf("Messages: %s\n", sessionIDReceiver)

	//get the message 

	messageData := pgsqlcrud.GetMessageData(ctx, config.MessageTable, pool,  data.ChatID , sender)

	// print the message data
	fmt.Printf("Messages: %+v\n", messageData)


}
