/* Make the redis connector 

*/

package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// ConnectRedis creates and returns a Redis client
func ConnectRedis(host string, port int) (*redis.Client, error) {
	
	addr := fmt.Sprintf("%s:%d", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // No password
		DB:       0,  // Use default DB
		Protocol: 2,  // RESP2
	})

	// Optionally test the connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	fmt.Println("Connected to Redis!")
	return client, nil
}