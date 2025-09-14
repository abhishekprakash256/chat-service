package ws


import (

	"net/http"

	"chat-service/internal/chat/wsendpoint"
	//"chat-service/internal/config"
)





func WsHandler(mux *http.ServeMux) {


    mux.HandleFunc("/chat-server/ws", wsendpoint.WSEndpoint) 
       
	
}


