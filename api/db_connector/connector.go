
package db_connector


import (
	"log"

	pgsqlconn "chat-service/internal/storage/pgsql/connection"

	redisconn "chat-service/internal/storage/redis/connection"

	"chat-service/internal/config"


)



func DbConnector() {

	client, errRedis := redisconn.ConnectRedis(config.RedisDefaultConfig.Host, config.RedisDefaultConfig.Port)


	if errRedis != nil {

		log.Fatalf("Failed to connect to Redis: %v", errRedis)

	}

		// Create the connection pool
	pool, errPgsql := pgsqlconn.ConnectPgSql(
		config.PgsqlDefaultConfig.Host,
		config.PgsqlDefaultConfig.User,
		config.PgsqlDefaultConfig.Password,
		config.PgsqlDefaultConfig.DBName,
		config.PgsqlDefaultConfig.Port,
	)

	// Create the database schema
	// The connection failed
	if errPgsql != nil {
		log.Fatal("DB connection failed: %v", errPgsql)

	}

	config.GlobalDbConn = &config.DbConn{
		PgsqlConn : pool , 
		RedisConn : client ,
	}

	log.Println("DB and Redis connection initialized successfully")


}

