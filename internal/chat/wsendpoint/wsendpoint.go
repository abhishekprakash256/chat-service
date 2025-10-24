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




// When a new connection arrives
// to add multiple connections in one sessionID
func AddClient(wsKey string , sessionKey string , conn *websocket.Conn) {

    config.ClientsWsMapper.Lock()
    defer config.ClientsWsMapper.Unlock()

    if config.ClientsWsMapper.Data[wsKey] == nil {
        config.ClientsWsMapper.Data[wsKey] = make(map[string]*websocket.Conn)
    }

    config.ClientsWsMapper.Data[wsKey][sessionKey] = conn
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



/*
func WSEndpoint(w http.ResponseWriter, r *http.Request) {
    log.Printf("Incoming WS request: %s", r.URL.RawQuery)

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Upgrade failed: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Println("WebSocket connection established")

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Disconnected:", err)
            break
        }
        log.Println("Received:", string(msg))
    }
}
*/



func WSEndpoint(w http.ResponseWriter, r *http.Request) {

	// get the data from url 
    chatID := r.URL.Query().Get("chatID")
	sessionID := r.URL.Query().Get("sessionID")
    sender := r.URL.Query().Get("user")

	log.Printf("chatid is %s", chatID)

	log.Printf("sessionID is %s", sessionID)

	log.Printf("sender is %s", sender)


	// get the session id as per machine 
	log.Printf("WS connection started")

    if chatID == "" || sender == ""  || sessionID == ""{ 
		log.Printf("WS connection started2")
        http.Error(w, "Missing chatID or user or sessionID ", http.StatusBadRequest)
		log.Printf("Some value is missing")
		log.Printf("%s,%s,%s", chatID , sessionID , sender)
        return
    }

	log.Printf("WS connection started3")

	// Step 1: Validate against DB
	ctx := context.Background()
	pool := config.GlobalDbConn.PgsqlConn

	log.Printf("WS connection started4")

	loginData, err := pgsqlcrud.GetLoginData(ctx, config.LoginTable , pool, chatID)

	log.Printf("WS connection started5")

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
	AddClient(wsKey , sessionKey , conn )

    log.Printf("Client connected: %s", sessionKey)

    // Step 4: Start session + heartbeat
	session.StartSession(conn , chatID , sessionID , sender )

	//go start the messagehub
	go messagehub.HandleMessages()

	messagereader.ReadMessage(conn)


}


