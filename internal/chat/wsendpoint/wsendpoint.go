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
//   - chatID  : chat session identifier
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
    chatID := r.URL.Query().Get("chatId")
	sessionID := r.URL.Query().Get("sessionID")
    sender := r.URL.Query().Get("user")

	// get the session id as per machine 


    if chatID == "" || sender == ""  || sessionID == ""{ 
        http.Error(w, "Missing chatID or user or sessionID ", http.StatusBadRequest)
        return
    }

	// Step 1: Validate against DB
	ctx := context.Background()
	pool := config.GlobalDbConn.PgsqlConn
	loginData, err := pgsqlcrud.GetLoginData(ctx, config.LoginTable , pool, chatID)

	if err != nil {
		http.Error(w, "Invalid chatID", http.StatusUnauthorized)
		log.Printf("WS connection rejected: invalid chatID %s", chatID)
		return
	}

	if sender != loginData.UserOne && sender != loginData.UserTwo {

	http.Error(w, "Invalid user for this session", http.StatusUnauthorized)

	log.Printf("WS connection rejected: user %s not part of chat %s", sender, chatID)
	
	return
	
	}


	// Step 2: Upgrade to WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WS upgrade failed: %v", err)
        return
    }

	// make the session key 
	wsKey := fmt.Sprintf("session:%s:%s", chatID, sender)

	//make the sssion unique per ws connection
    sessionKey := fmt.Sprintf("session:%s:%s:%s", chatID, sender , sessionID)

	// save the session in global ws mapper
	AddClient(wsKey , sessionKey , conn *websocket.Conn)

    log.Printf("Client connected: %s", sessionID)

    // Step 4: Start session + heartbeat
	session.StartSession(conn , chatID , sessionID , sender )

	//go start the messagehub
	go messagehub.HandleMessages()

	messagereader.ReadMessage(conn)


}



// When a new connection arrives
// to add multiple connections in one sessionID
func AddClient(wsKey string , sessionKey string , conn *websocket.Conn) {

    ClientsWsMapper.Lock()
    defer ClientsWsMapper.Unlock()

    if ClientsWsMapper.Data[wsKey] == nil {
        ClientsWsMapper.Data[wsKey] = make(map[string]*websocket.Conn)
    }

    ClientsWsMapper.Data[key][sessionKey] = conn
}
