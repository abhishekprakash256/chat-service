/*
the main to test the function.
*/

package main

import (
	"chat-service/internal/config"
	"chat-service/internal/hash"
	"fmt"
	"log"

	pgsqlconn "github.com/abhishekprakash256/go-pgsql-helper-kit/pgsql/db/connection"
	redisconn "github.com/abhishekprakash256/go-redis-helper-kit/redis/db/connection"
)

func main() {

	genrated_hash := hash.GenerateRandomHash(5, 10)

	fmt.Println(genrated_hash)

	//ctx := context.Background()

	// Making the connection
	client, errRedis := redisconn.ConnectRedis(config.RedisDefaultConfig.Host, config.RedisDefaultConfig.Port)

	if errRedis != nil {

		log.Fatalf("Failed to connect to Redis: %v", errRedis)

	}

	defer client.Close()

	hash.GenerateUniqueHash(config.UniqueHashSet, config.UsedHashSet, 5, 10, 20, client)

	// pop the hash from the primary set and get the hash

	i := 0

	for i < 10 {

		uniqueHash := hash.PopUniqueHash(config.UniqueHashSet, config.UsedHashSet, client)

		fmt.Println(uniqueHash)

		i++

	}

	// Create the connection pool
	pool, err := pgsqlconn.ConnectPgSql(
		config.PgsqlDefaultConfig.Host,
		config.PgsqlDefaultConfig.User,
		config.PgsqlDefaultConfig.Password,
		config.PgsqlDefaultConfig.DBName,
		config.PgsqlDefaultConfig.Port,
	)

	defer pool.Close() // Ensures pool is closed when program exits

	// Create the database schema
	// The connection failed
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

}
