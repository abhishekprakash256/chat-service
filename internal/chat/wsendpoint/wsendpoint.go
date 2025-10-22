/*
The wsendpoint to make the chat-end point for the data. The end point is used for the chat and sending the data for the user and get the data. 
The data is fetched for the data end point. 
*/
package wsendpoint



import (

	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	"context"

	"github.com/gorilla/websocket"
	"chat-service/internal/config"
	"chat-service/internal/chat/session"
	pgsqlcrud "chat-service/internal/storage/pgsql/crud"
	"chat-service/internal/chat/messagereader"
	"chat-service/internal/chat/messagehub"


)


// make the upgrader
var upgrader = websocket.Upgrader {

	ReadBufferSize : 4128 ,
	WriteBufferSize : 4128 , 
	CheckOrigin : func(r *http.Request) bool { return true } , 

}


// WSEndpoint upgrades HTTP connection to WebSocket
// and validates that the user is part of the chat session.
//
// URL Params (query):
//   - hash  : chat session identifier
//   - user  : username attempting to connect
//
// Flow:
//   1. Validate query params
//   2. Confirm user belongs to the chat in Postgres
//   3. Upgrade connection to WS
//   4. Save WS connection in ClientsWsMapper
//   5. Start heartbeat + Redis session tracking

func WSEndpoint(w http.ResponseWriter, r *http.Request) {

	// get the data from url 
    hash := r.URL.Query().Get("hash")
    sender := r.URL.Query().Get("user")

	// get the session id as per machine 


    if hash == "" || sender == "" {
        http.Error(w, "Missing hash or user", http.StatusBadRequest)
        return
    }

	// Step 1: Validate against DB
	ctx := context.Background()
	pool := config.GlobalDbConn.PgsqlConn
	loginData, err := pgsqlcrud.GetLoginData(ctx, config.LoginTable , pool, hash)

	if err != nil {
		http.Error(w, "Invalid hash", http.StatusUnauthorized)
		log.Printf("WS connection rejected: invalid hash %s", hash)
		return
	}

	if sender != loginData.UserOne && sender != loginData.UserTwo {

	http.Error(w, "Invalid user for this session", http.StatusUnauthorized)

	log.Printf("WS connection rejected: user %s not part of chat %s", sender, hash)
	
	return
	
	}


	// Step 2: Upgrade to WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WS upgrade failed: %v", err)
        return
    }

	// make the session id , //make the sssion unique per ws connection
    sessionID := fmt.Sprintf("session:%s:%s", hash, sender)


	// save the session in global ws mapper
    //config.ClientsWsMapper[sessionID] = conn
	AddClient(sessionID , conn)

    log.Printf("Client connected: %s", sessionID)

    // Step 4: Start session + heartbeat
	session.StartSession(conn , hash , sender)

	//go start the messagehub
	go messagehub.HandleMessages()

	messagereader.ReadMessage(conn)


}



// When a new connection arrives
// to add multiple connections in one sessionID

// change this to nested dict for message delivery 
func AddClient(sessionID string, conn *websocket.Conn) {
    if _, ok := config.ClientsWsMapper[sessionID]; !ok {
        config.ClientsWsMapper[sessionID] = []*websocket.Conn{}
    }
    config.ClientsWsMapper[sessionID] = append(config.ClientsWsMapper[sessionID], conn)
}