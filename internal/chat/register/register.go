package register



import (
	"encoding/json"
	"net/http"
	"fmt"
	"context"
	"log"

	"chat-service/internal/config"
	"chat-service/internal/hash"


	pgsqlcrud "chat-service/internal/storage/pgsql/crud"

)






// The Registartion handler to register the page
// Take two input 
// Generate the hash
// save the hash and the userOne and userTwo in the login table in the pgsql table 
//return the hash in the front-end 
func UserRegistration(w http.ResponseWriter, r *http.Request) {
	
	// import the connector from config
	pool := config.GlobalDbConn.PgsqlConn
	
	client := config.GlobalDbConn.RedisConn

	ctx := context.Background()
	
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintln(w, "Registration handler hit!")

	// Decode the JSON body
	var data config.RegistrationtData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// generate the hash
	hash.GenerateUniqueHash(config.UniqueHashSet, config.UsedHashSet, 5, 10, 20, client)

	// get the unique hash
	uniqueHash := hash.PopUniqueHash(config.UniqueHashSet, config.UsedHashSet, client)

	//save the login data into the db
	logindata := config.LoginData{
		ChatID:  uniqueHash,
		UserOne: data.UserOne,
		UserTwo: data.UserTwo,
	}

	if !pgsqlcrud.InsertLoginData(ctx, "login", pool, logindata) {
		log.Println("Insert into login failed")
	}
	

	// Print values (in logs and as response)
	fmt.Printf("userOne: %s, userTwo: %s\n", data.UserOne, data.UserTwo)
	
	fmt.Fprintf(w, "UserOne: %s, UserTwo: %s , Hash : %s \n", data.UserOne, data.UserTwo , uniqueHash )
	

}
