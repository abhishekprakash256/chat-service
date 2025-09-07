/*
the http end points for the registration , login , logout and end chat 

registration -- > 

--> end point ---> /chat-server/register 

take args --> userOne , userTwo

--> generate the hash 

pass into the end point with hash

input json --> 

data : {

userHash: "abc123"
userOne : "Abhi"
userTwo : "Anny"

}

return 

Response Json -->

{
  "status": "OK",
  "code": 200,
  "hash": "abc123"
}



*/

package http

import (

	"net/http"

	"chat-service/internal/chat/register"
	"chat-service/internal/chat/login"
)




func RegistrationHandler() {

	http.HandleFunc("/chat-server/register", register.UserRegistration)

}


func LoginHander() {

	http.HandleFunc("/chat-server/login", login.LoginUser)
	
}



