/*
make the logout 
*/

package logout 


import (

	"fmt"
	"context"
	"net/http"
	"chat-service/internal/config"
)


func LogOutUser(w http.ResponseWriter, r *http.Request ) {

	// get the method
	if r.Method != http.MethodPost {

		writeError(w, http.StatusMethodNotAllowed, "Only POST allowed")
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
	ctx := context.Background()


	fmt.Println(" ....... logout clicked ...... ")

}