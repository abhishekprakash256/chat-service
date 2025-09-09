
package wsendpoint



import (

	//"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"chat-service/internal/config"

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
func WSEndpoint(sessionID string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("The client could not be connected: %v", err)
		return
	}

	// TODO: You might want a mutex around this if many clients connect concurrently
	config.ClientsWsMapper[sessionID] = conn

	log.Printf("Client connected for session: %s", sessionID)


	// Print all active connections  to test
	fmt.Println("=== Active WebSocket Connections ===")
	for key, c := range config.ClientsWsMapper {
		fmt.Printf("SessionID: %s, Conn: %p\n", key, c)
	}
	fmt.Println("===================================")
}