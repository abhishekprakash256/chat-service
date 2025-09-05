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
	"chat-service/redis_service"
	

)


func main() {

	hash := hash_generation_service.GenerateRandomHash(5,10)

	fmt.Println(hash)

	//ctx := context.Background()

	// Making the connection
	client, err := connection.ConnectRedis(redis_config.DefaultConfig.Host, redis_config.DefaultConfig.Port)

	if err != nil {

		log.Fatalf("Failed to connect to Redis: %v", err)

	}

	defer client.Close()

	redis_hash_service.GenerateUniqueHash(redis_config.UniqueHashSet , redis_config.UsedHashSet , 5,10 , 20 , client )

	// pop the hash from the primary set and get the hash 

	var i int 
	
	i = 0 

	for i < 10 {
	uniqueHash := redis_hash_service.PopUniqueHash(redis_config.UniqueHashSet , redis_config.UsedHashSet , client )


	fmt.Println(uniqueHash)
	
	i++ 

	}
}
