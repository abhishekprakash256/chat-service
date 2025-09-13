/*
make the logout for the client
*/

package logout 


import (

	"fmt"
	"encoding/json"
	"strings"
	"context"
	"net/http"
	"log"
	"chat-service/internal/config"
	rediscrud "chat-service/internal/storage/redis/crud"



)


type IncomingMessage struct ()


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
	sessionID := r.URL.Query().Get("session")

	// Or from JSON body
	if sessionID == "" {
		var body struct {
			Session string `json:"session"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err == nil {

			sessionID = body.Session
		}
	}

	//get the chat id
	if sessionID == "" || !strings.HasPrefix(sessionID, "session:") {

		writeError(w, http.StatusBadRequest, "Missing or invalid session ID")

		return
	}

	// make the redis client
	client := config.GlobalDbConn.RedisConn

	//context passing baground
	ctx := context.Background()

	// Delete session from Redis
	err := rediscrud.DeleteSessionData(ctx, client, sessionID)

	if err != nil {
		log.Printf("Failed to delete session %s: %v", sessionID, err)
		writeError(w, http.StatusInternalServerError, "Logout failed")
		return
	}

	// err passed
	if err != nil {

		log.Printf("Failed to delete session %s: %v", sessionID, err)
		writeError(w, http.StatusInternalServerError, "Logout failed")
		return
	}

	// Also close WebSocket if still in memory 
	if conn, ok := config.ClientsWsMapper[sessionID]; ok {

		conn.Close()

		delete(config.ClientsWsMapper, sessionID)

	}

	// Send success response
	resp := SuccessResponse{
		Status: "success",
		Code:   http.StatusOK,
		Message: fmt.Sprintf("Logged out and session %s removed", sessionID),
	}

	//make the header
	w.Header().Set("Content-Type", "application/json")
	
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(resp)

	log.Printf("User logged out, session %s removed", sessionID)

}