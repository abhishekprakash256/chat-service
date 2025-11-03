package config


import (

	"github.com/gorilla/websocket"
	"sync" 
	"time"

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


type IncomingMessage struct {

	ChatID   string `json:"chatid"`
    Sender   string `json:"sender"`
    Receiver string `json:"receiver"`
    Message  string `json:"message"` 
}



type OutgoingMessage struct {

	MessageID int64 `json:"messageid"`
	ChatID   string `json:"chatid"`
    Sender   string `json:"sender"`
    Receiver string `json:"receiver"`
    Message  string `json:"message"` 
	Timestamp time.Time  `json:"time"`
}

// make the global dictonary for session:hash:name  to ws clinet 

//var ClientsWsMapper = make(map[string]*websocket.Conn) 

// ClientsWsMapper holds all active WebSocket connections per session:user
// change to nested mapepr for chatID and sessionID storage
var ClientsWsMapper = struct {
    sync.RWMutex
    Data map[string]map[string]*websocket.Conn
}{
    Data: make(map[string]map[string]*websocket.Conn),
}


//broadcast channel for broadcasting message
var BroadCast = make(chan []byte)