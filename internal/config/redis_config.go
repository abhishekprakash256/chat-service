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

type RedisSessionData struct {
	Hash      string
	Sender       string
	Reciever	string
	LastSeen    time.Time
	WSConnected int
	Notify      int
}

var UsedHashSet = "used_hash_set"

var UniqueHashSet = "frest_hash_set"

//unique_hash_set := "frest_hash_set"
