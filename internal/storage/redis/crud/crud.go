/*
The redis crud operation file for doing the operaration in redis
*/





package redis


import (
	"context"
	"time"
	"strconv"
	"github.com/redis/go-redis/v9"
	"chat-service/internal/config"
	"fmt"
)



func StoreSessionData(ctx context.Context, rdb *redis.Client, key string, data config.RedisSessionData) bool {

	err:= rdb.HSet(ctx, key, map[string]interface{}{
		"chat_id":      data.Hash,
		"sender":         data.Sender,
		"reciever":       data.Reciever , 
		"last_seen":    data.LastSeen.Format(time.RFC3339),
		"ws_connected": data.WSConnected,
		"notify":       data.Notify,
	}).Err()

	if err != nil {
		fmt.Println("The data pushed failed" , err)
		return false
	}

	fmt.Println("Data Pushed Succesfully ")
	return true

}

func GetSessionData(ctx context.Context, rdb *redis.Client, key string) (config.RedisSessionData, error) {
	result, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return config.RedisSessionData{}, err
	}

	// Convert ws_connected and notify from string to int
	wsConnected, err := strconv.Atoi(result["ws_connected"])
	if err != nil {
		wsConnected = 0
	}

	notify, err := strconv.Atoi(result["notify"])
	if err != nil {
		notify = 0
	}

	lastSeen, err := time.Parse(time.RFC3339, result["last_seen"])
	if err != nil {
		lastSeen = time.Time{}
	}

	return config.RedisSessionData{
		Hash:      result["chat_id"],
		Sender:        result["sender"],
		Reciever:      result["reciever"],
		LastSeen:    lastSeen,
		WSConnected: wsConnected,
		Notify:      notify,
	}, nil
}



func DeleteSessionData(ctx context.Context, rdb *redis.Client, key string) bool {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		fmt.Println("Failed to delete session data:", err)
		return false
	}
	return true
}

