/*
Mkae the login and start the session 


test data 


for registration -- 

{
    "userOne": "Sam",
     "userTwo": "Bob"
}

{
    "status": "success",
    "code": 200,
    "data": {
        "hash": "enyVF5JkoV0",
        "userOne": "Sam",
        "userTwo": "Bob"
    }
}


The key and the field -- 

| Redis Key             | Data Structure | Fields (example)                                                         |
| --------------------- | -------------- | ------------------------------------------------------------------------ |
| `session:abc123:Abhi` | Hash           | `ws_connected: 1` <br> `last_seen: 2025-07-08T20:00:00` <br> `notify: 0` |
| `session:abc123:Anny` | Hash           | `ws_connected: 0` <br> `last_seen: 2025-07-08T19:55:00` <br> `notify: 1` |
| `session:def456:Bob`  | Hash           | `ws_connected: 1` <br> `last_seen: 2025-07-08T20:01:00` <br> `notify: 0` |
| `session:def456:Cara` | Hash           | `ws_connected: 0` <br> `last_seen: 2025-07-08T19:50:00` <br> `notify: 1` |


HSET session:abc123:Abhi chat_id abc123
HSET session:abc123:Abhi user Abhi
HSET session:abc123:Abhi last_seen 2025-07-08T20:00:00
HSET session:abc123:Abhi ws_connected 1
HSET session:abc123:Abhi notify 0

when login store the session data into the redis
and start hearbeat protocol until logout or endchat

*/

/*



func LoginUser 

params --> hash , user 

hash comes from the link

match the user in the pgsql
create the login
store the session data into the redis 
and start the hearbeat protocol to check the session 
and keep updating the status
pass a success message if user is valid to front-end 

get the json 

{

UserName : "Abhi"
hash : "abc123"

}

returns 

{
    status : OK
    code : 200
}


type ErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}


*/

package login 


import (

	"net/http"
	"encoding/json"
	"fmt"
	"context"
	"log"


	"chat-service/internal/config"

	pgsqlcrud "chat-service/internal/storage/pgsql/crud"

)

type LoginRequest struct {
    Hash     string `json:"hash"`
    UserName string `json:"username"`
}


type LoginSuccess struct {
	
	Status  string `json:"status"`
	Code    int    `json:"code"`
	

}


type ErrorResponse struct {
	
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`

}


func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  "error",
		Code:    code,
		Message: msg,
	})
}



func LoginUser( w http.ResponseWriter, r *http.Request ) {

	ctx := context.Background()

	// make the redis and pgsql connection
	pool := config.GlobalDbConn.PgsqlConn

	//client := config.GlobalDbConn.RedisConn. not used rn 

	if r.Method != http.MethodPost {
        writeError(w, http.StatusMethodNotAllowed, "Only POST allowed")
        return
    }

	// Decode request
	var data LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// for testing purpose
	fmt.Printf("Login attempt: %s for chat %s\n", data.UserName, data.Hash)

	//get the login 
	retrievedLogin, err := pgsqlcrud.GetLoginData(ctx, "login", pool,  data.Hash )

	// get the login details
	if err != nil {
		log.Println("Login not found:", err)
	
		} else {
		fmt.Printf("Login for chat %s: %s & %s\n", retrievedLogin.ChatID, retrievedLogin.UserOne, retrievedLogin.UserTwo)
	}

	//match the login data 
	if data.UserName == retrievedLogin.UserOne  || data.UserName == retrievedLogin.UserTwo {

		log.Printf("Login succefull for %s" , data.UserName)

	} else {
		
		log.Printf("failed login for %s" , data.UserName)

		// return the fail login data 


	}

	//make the json 



	// start the session 


	// return the succesfll json


}



