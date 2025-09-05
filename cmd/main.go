/*
the main to test the function.
*/

package main

import (
	"chat-service/internal/config"
	"chat-service/internal/hash"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abhishekprakash256/go-pgsql-helper-kit/pgsql/db/crud"

	pgsqlconn "github.com/abhishekprakash256/go-pgsql-helper-kit/pgsql/db/connection"
	redisconn "github.com/abhishekprakash256/go-redis-helper-kit/redis/db/connection"
)

func main() {

	genrated_hash := hash.GenerateRandomHash(5, 10)

	fmt.Println(genrated_hash)

	ctx := context.Background()

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

	// Create the database schema
	// The connection failed
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	// Create the database schema
	if err := crud.CreateSchema(ctx, pool, config.LoginTableSQL, config.MessageTableSQL); err != nil {

		log.Fatal("Schema creation failed:", err)
	}

	// Test login data
	login := crud.LoginData{
		ChatID:  "abc123",
		UserOne: "Abhi",
		UserTwo: "Anny",
	}

	// login the data
	if !crud.InsertLoginData(ctx, "login", pool, login) {
		log.Println("Insert into login failed")
	}

	// Test message data
	msg := crud.MessageData{
		ChatID:       "abc123",
		SenderName:   "Abhi",
		ReceiverName: "Anny",
		Message:      "Hello There!",
		Timestamp:    time.Now(),
		Read:         false,
	}

	// Insert the message data
	if !crud.InsertMessageData(ctx, "message", pool, msg) {
		log.Println("Insert into message failed")
	}

	// Step 5: Retrieve login data
	retrievedLogin, err := crud.GetLoginData(ctx, "login", pool, "abc123")
	if err != nil {
		log.Println("Login not found:", err)
	} else {
		fmt.Printf("Login for chat %s: %s & %s\n", retrievedLogin.ChatID, retrievedLogin.UserOne, retrievedLogin.UserTwo)
	}

	// Step 6: Retrieve message data
	messages := crud.GetMessageData(ctx, "message", pool, "abc123", "Abhi")

	// print the message data
	fmt.Printf("Messages: %+v\n", messages)

	// to print the message from the database
	for _, m := range messages {
		fmt.Printf("Message from %s to %s: %s\n", m.SenderName, m.ReceiverName, m.Message, m.Timestamp, m.Read)
	}

	// Test delete message data
	if !crud.DeleteMessageData(ctx, "message", pool, "abc123") {
		log.Println("Delete message failed")
	}

	// Test delete login data
	if !crud.DeleteLoginData(ctx, "login", pool, "abc123") {
		log.Println("Delete login data failed")
	}

	//Done with the operation
	log.Println("Data operation done successfully")

	defer pool.Close() // Ensures pool is closed when program exits
}
