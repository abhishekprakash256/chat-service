package config


import (

	"github.com/gorilla/websocket"

	pgsqlconn "github.com/jackc/pgx/v5/pgxpool"

	redisconn "github.com/redis/go-redis/v9"
)


type DbConn struct {

	PgsqlConn *pgsqlconn.Pool

	RedisConn *redisconn.Client
}

var GlobalDbConn *DbConn


// Define the struct to match the expected JSON
type RegistrationtData struct {
	UserOne string `json:"userOne"`
	UserTwo string `json:"userTwo"`
}



// make the global dictonary for session:hash:name  to ws clinet 

//var ClientsWsMapper = make(map[string]*websocket.Conn) 

// ClientsWsMapper holds all active WebSocket connections per session:user
var ClientsWsMapper = make(map[string][]*websocket.Conn)

//broadcast channel for testing 
var BroadCast = make(chan []byte)