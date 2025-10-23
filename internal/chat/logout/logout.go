/*
make the logout for the client
*/

package logout 


import (

	"fmt"
	"encoding/json"
	"time"
	"context"
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"chat-service/internal/config"
	"chat-service/internal/chat/session"
	pgsqlcrud "chat-service/internal/storage/pgsql/crud"
	//rediscrud "chat-service/internal/storage/redis/crud"

	
)


type LogoutRequest struct  {

	ChatID string `json:"ChatID"`
	SessionID string `json:"SessionID"`
	UserName string `json:"UserName"`
}


// SuccessResponse defines a standard logout success payload
type SuccessResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse defines a standard error payload
type ErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}





// helper to write error consistently
func writeError(w http.ResponseWriter, code int, msg string) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  "error",
		Code:    code,
		Message: msg,
	})

}

// the logout function will take a post request 
// get the sessionid  which has chatid and userid 
// find the sesisonid in redis string and then upodate the session data
// stop the session 
// delete the currrent ws client connection from the clinetwsmapper
func LogOutUser(w http.ResponseWriter, r *http.Request ) {

	// get the method
	if r.Method != http.MethodPost {

		writeError( w , http.StatusMethodNotAllowed, "Only POST allowed" )
		return

	}

	// Try to get sessionID from query param: /logout?session=session:abc:User1
	// Decode request JSON
	var data LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Debug log
	fmt.Printf( "Logout attempt: %s for chat %s\n", data.UserName, data.ChatID )

	// get the pool 
	pool := config.GlobalDbConn.PgsqlConn

	//context passing baground
	ctx := context.Background()

	// Retrieve login record from DB
	retrievedLogin, err := pgsqlcrud.GetLoginData(ctx, config.LoginTable , pool, data.ChatID)

	if err != nil {
		
		writeError(w, http.StatusUnauthorized, "Login Failed: Hash not found")
		log.Printf("Login failed: hash %s not found (%v)", data.ChatID, err)
		return
	}

	// make the sender and reciver 
	var sender string 
	var receiver string 

	// Check if username matches registered users
	if data.UserName == retrievedLogin.UserOne {
	
		sender = retrievedLogin.UserOne
		receiver = retrievedLogin.UserTwo
	} else if data.UserName == retrievedLogin.UserTwo {
		
		sender = retrievedLogin.UserTwo
		receiver = retrievedLogin.UserOne
	} else {
		
		// username does not match either registered user → fail login
		writeError(w, http.StatusUnauthorized, "Login Failed: Wrong Username or Hash")
		log.Printf("Login failed: username %s not valid for hash %s", data.UserName, data.ChatID)
		return
	}


	// data to save in the redis sesssion
	now := time.Now()
	ws_connected := 0
	notify := 0 

	// save the data with ws connected 1 
	session.SaveSession(data.ChatID, data.SessionID , sender , receiver , now , ws_connected, notify)


	wsKey := fmt.Sprintf("session:%s:%s", data.ChatID, sender)
	// make the sessionid 
	sessionKey := fmt.Sprintf("session:%s:%s:%s", data.ChatID, sender, data.SessionID)


	//go through all the connection and delete all 
	// --- Thread-safe WS cleanup ---
	config.ClientsWsMapper.Lock()

	if sessionMap, ok := config.ClientsWsMapper.Data[wsKey]; ok {
		if conn, exists := sessionMap[sessionKey]; exists {
			_ = conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Logged out"),
			)
			conn.Close()
			delete(sessionMap, sessionKey)
			log.Printf("Closed and removed WebSocket for %s", sessionKey)
		}

		// Clean up top-level key if no sessions remain
		if len(sessionMap) == 0 {
			delete(config.ClientsWsMapper.Data, wsKey)
		}
	}

	config.ClientsWsMapper.Unlock()

	// Send success response
	resp := SuccessResponse{
		Status: "success",
		Code:   http.StatusOK,
		Message: fmt.Sprintf("Logged out and session %s removed", sessionKey),
	}

	//make the header
	w.Header().Set("Content-Type", "application/json")
	
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(resp)

	log.Printf("User logged out, session %s removed", sessionKey)

	
	}
