/*
the main to test the function. 
*/


package main 

import (

	"chat-service/chat_hash_service"
	"log"
	"fmt"
	"github.com/abhishekprakash256/go-redis-helper-kit/redis/db/connection"
	"chat-service/redis_service"

)


func main() {

	num := chat_hash_generation.GenerateRandomHash(5,10)

	fmt.Println(num)

	//ctx := context.Background()

	// Making the connection
	client, err := connection.ConnectRedis(config.DefaultConfig.Host, config.DefaultConfig.Port)

	if err != nil {

		log.Fatalf("Failed to connect to Redis: %v", err)

	}

	defer client.Close()
	
}