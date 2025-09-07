package config


import (
	pgsqlconn "chat-service/internal/storage/pgsql/connection"

	redisconn "chat-service/internal/storage/redis/connection"
)


type DbConn struct {

	PgsqlConn *pgsqlconn.ConnectRedis

	RedisConn *redisconn.ConnectPgSql
}

var GlobalDbConn *DbConn


// Define the struct to match the expected JSON
type RegistrationtData struct {
	UserOne string `json:"userOne"`
	UserTwo string `json:"userTwo"`
}



