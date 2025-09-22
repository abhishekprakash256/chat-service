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
	"chat-service/internal/chat/logout"
	"chat-service/internal/chat/endchat"
)



/*
func RegistrationHandler() {

	http.HandleFunc("/chat-server/user/register", register.UserRegistration)

}


func LoginHander() {

	http.HandleFunc("/chat-server/user/login", login.LoginUser)

}

func LogoutHandler() {

	http.HandleFunc("/chat-server/user/logout", logout.LogOutUser)

}

*/



// SetupUserRoutes attaches user routes to the given mux.
func SetupUserRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/chat-server/user/register", register.UserRegistration)
    mux.HandleFunc("/chat-server/user/login", login.LoginUser)
    mux.HandleFunc("/chat-server/user/logout", logout.LogOutUser)
	mux.HandleFunc("/chat-server/user/endchat" , endchat.UserEndChat)
}



