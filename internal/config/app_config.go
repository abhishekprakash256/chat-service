package config


import (
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



