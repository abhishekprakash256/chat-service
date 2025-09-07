package register



import (
	"encoding/json"
	"net/http"
	"fmt"
	"chat-service/internal/config"
)






// The Registartion handler to register the page
// Take two input 
// Generate the hash
// save the hash and the userOne and userTwo in the login table in the pgsql table 
//return the hash in the front-end 
func UserRegistration(w http.ResponseWriter, r *http.Request) {
	
	
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

	// Print values (in logs and as response)
	fmt.Printf("userOne: %s, userTwo: %s\n", data.UserOne, data.UserTwo)
	fmt.Fprintf(w, "Received -> UserOne: %s, UserTwo: %s\n", data.UserOne, data.UserTwo)
	

}
