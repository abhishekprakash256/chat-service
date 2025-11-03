/*

The redis config for the  connection
*/

package config

import (

	"time"
)

type RedisDBConfig struct {
	Host string
	Port int
}

var RedisDefaultConfig = RedisDBConfig{
	Host: "localhost",
	Port: 6379,
}


//change this as per SessionID as well and make hash as ChatID all places
//chnage Hash to ChatID and 
type RedisSessionData struct {
	ChatID      string
	SessionID   string
	Sender       string
	Reciever	string
	LastSeen    time.Time
	WSConnected int
	Notify      int
}



var UsedHashSet = "used_hash_set"

var UniqueHashSet = "frest_hash_set"

//unique_hash_set := "frest_hash_set"
