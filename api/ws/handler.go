package ws


import (

	"net/http"

	"chat-service/internal/chat/wsendpoint"
	//"chat-service/internal/config"
)





func WsHandler(sessionID string) {


    http.HandleFunc("/chat-server/ws", func(w http.ResponseWriter, r *http.Request) {
        wsendpoint.WSEndpoint(sessionID, w, r)
    })
	
}
