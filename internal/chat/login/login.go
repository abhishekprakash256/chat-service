/*
Mkae the login and start the session 


test data 


for registration  Testing-- 

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

{
    "status": "success",
    "code": 200,
    "data": {
        "hash": "XWVU7wbbr",
        "userOne": "Bob",
        "userTwo": "Ben"
    }
}


////////-----------------------------------////////

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
    "status": "success",
    "code": 200,
    "data": {
        "hash": "enyVF5JkoV0",
        "sender": "Sam",
		"reciever": "Ben"
    }	
}


type ErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}


*/


package login



import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"chat-service/internal/chat/session"
	"chat-service/internal/config"
	//wshandler "chat-service/api/ws"
	pgsqlcrud "chat-service/internal/storage/pgsql/crud"
)

// LoginRequest represents the expected JSON body for login requests.
type LoginRequest struct {
	Hash     string `json:"Hash"`     // Chat session hash
	UserName string `json:"UserName"` // Username trying to log in
}

// LoginSuccess represents a successful login response.
type LoginSuccess struct {
	Status string      `json:"status"` // Always "success"
	Code   int         `json:"code"`   // HTTP status code
	Data   interface{} `json:"data"`   // Additional data payload
}

// ErrorResponse represents a standardized error response.
type ErrorResponse struct {
	Status  string `json:"status"`  // Always "error"
	Code    int    `json:"code"`    // HTTP status code
	Message string `json:"message"` // Human-readable error message
}

// writeError writes a standardized JSON error response.
func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  "error",
		Code:    code,
		Message: msg,
	})
}




// LoginUser handles chat login requests.
//
// Expected request (POST /chat-server/login):
//   {
//     "hash": "abc123",
//     "username": "Sam"
//   }
//
// Steps performed:
//   1. Validates that the request method is POST.
//   2. Parses the JSON request body into LoginRequest.
//   3. Looks up login data by hash in PostgreSQL.
//   4. Checks if the username matches one of the registered users.
//   5. Starts a session in Redis (via session.StartSession).
//   6. Returns a standardized JSON response.
//
// Success Response Example:
//   {
//     "status": "success",
//     "code": 200,
//     "data": {
//       "hash": "abc123",
//       "username": "Sam"
//     }
//   }
//
// Error Response Example:
//   {
//     "status": "error",
//     "code": 401,
//     "message": "Login Failed Wrong Username or Hash"
//   }
// when login clicked the ChatID and the username will be stored in the localsession storage for the front end 
// using the sessionStortage 
// as this in the json
// {  
//  "ChatID": "XWVU7wbbr",
//  "UserName": "Ben"
// }

func LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	pool := config.GlobalDbConn.PgsqlConn

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	// Decode request JSON
	var data LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Debug log
	fmt.Printf("Login attempt: %s for chat %s\n", data.UserName, data.Hash)

	// Retrieve login record from DB
	retrievedLogin, err := pgsqlcrud.GetLoginData(ctx, config.LoginTable , pool, data.Hash)
	
	if err != nil {
		
		writeError(w, http.StatusUnauthorized, "Login Failed: Hash not found")
		log.Printf("Login failed: hash %s not found (%v)", data.Hash, err)
		return
	}

	// make the sender and reciver 
	var sender string 
	var receiver string 

	// Check if username matches registered users
	if data.UserName == retrievedLogin.UserOne {
	
		sender = retrievedLogin.UserOne
		receiver = retrievedLogin.UserTwo
	} else if data.UserName == retrievedLogin.UserTwo {
		
		sender = retrievedLogin.UserTwo
		receiver = retrievedLogin.UserOne
	} else {
		
		// username does not match either registered user → fail login
		writeError(w, http.StatusUnauthorized, "Login Failed: Wrong Username or Hash")
		log.Printf("Login failed: username %s not valid for hash %s", data.UserName, data.Hash)
		return
	}

	// Successful login
	log.Printf("Login successful for %s (chat %s)", data.UserName, data.Hash)

	// StartSession(hash string , sender string , reciever string )
	log.Printf("Sender %s and Reciever %s", sender, receiver )

	// start the session
	//session.StartSession(data.Hash , sender ,  receiver )

	// data to save in the redis sesssion
	now := time.Now()
	ws_connected := 0
	notify := 0 

	// save the data
	session.SaveSession(data.Hash, sender , receiver , now , ws_connected, notify)

	// Success response
	resp := LoginSuccess{
		Status: "success",
		Code:   http.StatusOK,
		Data: map[string]string{
			"hash":     data.Hash,
			"sender": data.UserName,
			"receiver" : receiver , 
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}