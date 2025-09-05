/*
the main to test the function.
*/

package main

import (
	"chat-service/internal/config"
	"chat-service/internal/hash"
	"fmt"
	"log"

	redisconn "github.com/abhishekprakash256/go-redis-helper-kit/redis/db/connection"
	//pgsqlconn "github.com/abhishekprakash256/go-pgsql-helper-kit/pgsql/db/connection"
)

func main() {

	genrated_hash := hash.GenerateRandomHash(5, 10)

	fmt.Println(genrated_hash)

	//ctx := context.Background()

	// Making the connection
	client, err := redisconn.ConnectRedis(config.RedisDefaultConfig.Host, config.RedisDefaultConfig.Port)

	if err != nil {

		log.Fatalf("Failed to connect to Redis: %v", err)

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
}
