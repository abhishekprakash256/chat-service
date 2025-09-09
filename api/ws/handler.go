package ws


import (

	"net/http"

	"chat-service/internal/chat/wsendpoint"
	//"chat-service/internal/config"
)





func WsHandler() {


    http.HandleFunc("/chat-server/ws", wsendpoint.WSEndpoint) 
       
	
}


