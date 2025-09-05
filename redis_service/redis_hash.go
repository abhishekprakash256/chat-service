/*
The function to make the store the hash into the redis hash set 

*/

package redis_hash_service

import (

	"chat-service/hash_service"
	"log"
	"github.com/redis/go-redis/v9"
	"context"

)


// GenerateUniqueHash generates unique random hashes and stores them in Redis.
//
// Parameters:
//   - uniqueHashSet: Redis set name for fresh hashes.
//   - usedHashSet: Redis set name for used hashes (not used in this function yet).
//   - hostname: The server hostname (not used yet, but can be helpful for sharding).
//   - minHashSize: Minimum size of generated hash.
//   - maxHashSize: Maximum size of generated hash.
//   - hashQty: Number of unique hashes to generate.
//   - redisClient: Optional *redis.Client. If nil, a new client is created from config.
//
// Behavior:
//   1. Generate a random hash between minHashSize and maxHashSize.
//   2. Check if it already exists in the unique set.
//   3. If not, add it to the unique set.
//   4. Continue until `hashQty` hashes are stored.

func GenerateUniqueHash(
	uniqueHashSet string,
	usedHashSet string,
	minHashSize int,
	maxHashSize int,
	hashQty int,
	redisClient *redis.Client,
) {
	ctx := context.Background()
	client := redisClient

	for {
		// Step 0: Check how many hashes are already stored
		currentCount, err := client.SCard(ctx, uniqueHashSet).Result()
		if err != nil {
			log.Printf("Redis SCard error: %v", err)
			return
		}

		if currentCount >= int64(hashQty) {
			log.Printf("Already have %d hashes in set %s. Stopping.", currentCount, uniqueHashSet)
			break
		}

		// Step 1: Generate a random hash
		hash := hash_generation_service.GenerateRandomHash(minHashSize, maxHashSize)

		// Step 2: Check if it already exists
		exists, err := client.SIsMember(ctx, uniqueHashSet, hash).Result()
		if err != nil {
			log.Printf("Redis SIsMember error: %v", err)
			continue
		}

		if exists {
			continue // skip duplicates
		}

		// Step 3: Add new hash
		_, err = client.SAdd(ctx, uniqueHashSet, hash).Result()
		if err != nil {
			log.Printf("Error adding hash to Redis set: %v", err)
			continue
		}

		log.Printf("Generated new unique hash: %s", hash)
	}

	log.Printf("Finished ensuring %d unique hashes exist in set %s", hashQty, uniqueHashSet)
}




// PopUniqueHash pops a random hash from the `uniqueHashSet` in Redis
// and optionally adds it to `usedHashSet` to mark it as used.
//
// Parameters:
//   - uniqueHashSet: Redis set containing fresh hashes.
//   - usedHashSet: Redis set to store used hashes (can be empty string if not needed).
//   - redisClient: Redis client instance.
//
// Returns:
//   - hash string if successfully popped, empty string otherwise.
func PopUniqueHash(uniqueHashSet string, usedHashSet string, redisClient *redis.Client) string {

	ctx := context.Background()

	// Step 1: Pop a random element from the set
	hash, err := redisClient.SPop(ctx, uniqueHashSet).Result()
	if err == redis.Nil {

		log.Printf("No hashes left in set %s", uniqueHashSet)
		return ""

	} else if err != nil {

		log.Printf("Error popping hash from set: %v", err)

		return ""
	}

	// Step 2: Move the hash to usedHashSet

	_, err = redisClient.SAdd(ctx, usedHashSet, hash).Result()

	if err != nil {

		log.Printf("Error adding hash to used set: %v", err)
	}


	log.Printf("Popped hash: %s", hash)
	return hash

	
}