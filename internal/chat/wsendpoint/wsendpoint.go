
package wsendpoint



import (

	//"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"chat-service/internal/config"
	"chat-service/internal/chat/session"

)


// make the upgrader
var upgrader = websocket.Upgrader {

	ReadBufferSize : 4128 ,
	WriteBufferSize : 4128 , 
	CheckOrigin : func(r *http.Request) bool { return true } , 

}


// the ws end point to connect the user 
// route the message 
// save the message 
// check the status 
// update in redis as well 


// WSEndpoint upgrades HTTP connection to WebSocket and saves it in the mapper.
func WSEndpoint(w http.ResponseWriter, r *http.Request) {

	// get the data from url 
    hash := r.URL.Query().Get("hash")
    sender := r.URL.Query().Get("user")


    if hash == "" || sender == "" {
        http.Error(w, "Missing hash or user", http.StatusBadRequest)
        return
    }

	// make the session id 
    sessionID := fmt.Sprintf("session:%s:%s", hash, sender)

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WS upgrade failed: %v", err)
        return
    }


	// save the session in global ws mapper
    config.ClientsWsMapper[sessionID] = conn
    log.Printf("Client connected: %s", sessionID)

    //start session
	session.StartSession(conn , hash , sender)


}
