package ws


import (

	"net/http"

	"chat-service/internal/chat/wsendpoint"
	//"chat-service/internal/config"
)





func WsHandler(mux *http.ServeMux) {


    mux.HandleFunc("/chat-server/v1/ws/chat", wsendpoint.WSEndpoint) 
       
	
}


