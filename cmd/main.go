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
	"net/http"
	"github.com/rs/cors"

	"chat-service/api/db_connector"

	pgsqlcrud "chat-service/internal/storage/pgsql/crud"

	pgsqldb "chat-service/internal/storage/pgsql/db"

	httphandler "chat-service/api/http"

	wshandler "chat-service/api/ws"
)



/*
func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Registration handler hit!")
}
*/

func main() {

	// run the code to make the connetion
	db_connector.DbConnector()

	genrated_hash := hash.GenerateRandomHash(5, 10)

	fmt.Println(genrated_hash)

	ctx := context.Background()

	// import the connector from db
	pool := config.GlobalDbConn.PgsqlConn
	
	client := config.GlobalDbConn.RedisConn


	hash.GenerateUniqueHash(config.UniqueHashSet, config.UsedHashSet, 5, 10, 20, client)

	// pop the hash from the primary set and get the hash

	i := 0

	for i < 10 {

		uniqueHash := hash.PopUniqueHash(config.UniqueHashSet, config.UsedHashSet, client)

		fmt.Println(uniqueHash)

		i++

	}


	// Create the database schema
	if err := pgsqldb.CreateSchema(ctx, pool, config.LoginTableSQL, config.MessageTableSQL); err != nil {

		log.Fatal("Schema creation failed:", err)
	}

	// Test login data
	login := config.LoginData{
		ChatID:  "abc123",
		UserOne: "Abhi",
		UserTwo: "Anny",
	}

	// login the data
	if !pgsqlcrud.InsertLoginData(ctx, "login", pool, login) {
		log.Println("Insert into login failed")
	}

	// Test message data
	msg := config.MessageData{
		ChatID:       "abc123",
		Sender:   "Abhi",
		Receiver: "Anny",
		Message:      "Hello There!",
		Timestamp:    time.Now(),
		Read:         false,
	}

	// Insert the message data
	if !pgsqlcrud.InsertMessageData(ctx, "message", pool, msg) {
		log.Println("Insert into message failed")
	}

	// Step 5: Retrieve login data
	retrievedLogin, err := pgsqlcrud.GetLoginData(ctx, "login", pool, "abc123")
	if err != nil {
		log.Println("Login not found:", err)
	} else {
		fmt.Printf("Login for chat %s: %s & %s\n", retrievedLogin.ChatID, retrievedLogin.UserOne, retrievedLogin.UserTwo)
	}

	// Step 6: Retrieve message data
	messages := pgsqlcrud.GetMessageData(ctx, "message", pool, "abc123", "Abhi")

	// print the message data
	fmt.Printf("Messages: %+v\n", messages)

	// to print the message from the database
	for _, m := range messages {
		fmt.Printf("Message from %s to %s: %s\n", m.Sender, m.Receiver, m.Message, m.Timestamp, m.Read)
	}

	// Test delete message data
	if !pgsqlcrud.DeleteMessageData(ctx, "message", pool, "abc123") {
		log.Println("Delete message failed")
	}

	// Test delete login data
	if !pgsqlcrud.DeleteLoginData(ctx, "login", pool, "abc123") {
		log.Println("Delete login data failed")
	}

	//Done with the operation
	log.Println("Data operation done successfully")

	defer pool.Close() // Ensures pool is closed when program exits

	// make the mux for server 
	mux := http.NewServeMux()

	// call the routes
	httphandler.SetupUserRoutes(mux)

	// start the WS handler (no CORS needed here)
	wshandler.WsHandler(mux)

	// wrap HTTP mux with CORS middleware
	c := cors.New(cors.Options{

		// for prod use https://meabhi.me
		AllowedOrigins:   []string{"http://localhost:3000"},  // use for dev
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	fmt.Println("Server started on :8050")
	if err := http.ListenAndServe(":8050", handler); err != nil {
		log.Fatal(err)
	}


	
}
