/*
the main to test the function. 
*/


package main 

import (

	"chat-service/hash_service"
	"log"
	"fmt"
	"github.com/abhishekprakash256/go-redis-helper-kit/redis/db/connection"
	"chat-service/redis_config"
	"chat-service/redis_hash_store"

)


func main() {

	hash := hash_generation.GenerateRandomHash(5,10)

	fmt.Println(hash)

	//ctx := context.Background()

	// Making the connection
	client, err := connection.ConnectRedis(config.DefaultConfig.Host, config.DefaultConfig.Port)

	if err != nil {

		log.Fatalf("Failed to connect to Redis: %v", err)

	}

	defer client.Close()

	hash_manager.GenerateUniqueHash(config.UniqueHashSet , config.UsedHashSet , 5,10 , 20 , client )

	// pop the hash from the primary set and get the hash 

	var i int 
	
	i = 0 

	for i < 10 {
	uniqueHash := hash_manager.PopUniqueHash(config.UniqueHashSet , config.UsedHashSet , client )


	fmt.Println(uniqueHash)
	
	i++ 
	
	}
}
