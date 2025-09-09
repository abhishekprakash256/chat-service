/*
make the session


type RedisSessionData struct {
	ChatID      string
	Sender       string
	Reciever	string
	LastSeen    time.Time
	WSConnected int
	Notify      int
}

session:abc123:Abhi → {
  chat_id: abc123,
  sender: Abhi,
  reciver: Anny,
  last_seen: 2025-07-08T20:00:00,
  ws_connected: 1,
  notify: 0
}

*/

package session 

import (

	"fmt"
	"chat-service/api/db_connector"
	"time"

	"chat-service/internal/config"
	rediscrud "chat-service/internal/redis/crud"

)


// save the session
// params --> hash , sender , reciever ,time ,  ws_status , notify value  
// get the client 
// make the struct 
func SaveSession(hash string , sender string , reciever string , last_seen time.Time ,  ws_status int , notification int  ) {

	client := config.GlobalDbConn.RedisConn

	ctx := context.Background()

	sessionData := config.RedisSessionData{

		Hash : hash ,
		Sender : sender , 
		Reciever : reciever , 
		LastSeen : last_seen , 
		WSConnected : ws_status , 
		Notify : notification , 

	}

	var sessionId string 

	sessionId = "session":hash:Sender

	val := rediscrud.StoreSessionData(ctx , client , sessionId , sessionData)

	if val != nil {
		fmt.Println("Session data has been saved")

		return true
	}

	fmt.Println("Session data not saved")
	
	return false

}




// strat the session
// The save session take the redis client 
// params -- > usename , hash 
// make the redis hash , take  timestamp and save the data 

func StartSession(hash string , sender string , reciever string ) {

	fmt.Println("Session started")

	//client := config.GlobalDbConn.RedisConn

	//ctx := context.Background()

	SaveSession()
	
}



