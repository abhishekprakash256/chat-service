/*
The function to make the store the hash into the redis hash set 

*/

package generate_unique_hash

import (

	"chat-service/hash_service"
	"log"
	"fmt"
	"github.com/abhishekprakash256/go-redis-helper-kit/redis/db/connection"
	"chat-service/redis_config"
	"context"

)



func GenerateUniqueHash(
    uniqueHashSet string,
    usedHashSet string,
    hostname string,
    minHashSize int,
    maxHashSize int,
    hashQty int,
    redisClient *RedisClient, // pointer so it can be nil
	) {

	/*
	The function to generate the random hash
	compare the hash in the used hash set 
	if the hash is unique store it and pop out a random hash for use 
	*/

	client, err := connection.ConnectRedis(config.DefaultConfig.Host, config.DefaultConfig.Port)

	ctx := context.Background()
	
	// get the length of the unque hash set 
	cardinality, err := client.SCard(ctx, uniqueHashSet).Result()
	
	// iter to the length 
	i := 0
	
	for i < hashQty {

		hash := hash_generation.GenerateRandomHash(5,10)

		// check if the value not in the set 
		exists, err := client.SIsMember(ctx, uniqueHashSet, hash).Result()
		
		if err != nil {
			// Handle error
		}

		if exists {

			continue
		}

		else {

		_, err := rdb.SAdd(ctx, uniqueHashSet, hash).Result()

		if err != nil {
			log.Fatalf("Error adding value1 to set: %v", err)
		}

		}


		

		
	} 







}