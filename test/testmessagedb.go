/*
test message db for getting the message
*/

package main 



import (
	"chat-service/internal/config"

	"context"
	"fmt"
	//"log"

	"chat-service/api/db_connector"

	pgsqlcrud "chat-service/internal/storage/pgsql/crud"

)


/*
func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Registration handler hit!")
}
*/

func main() {

	// run the code to make the connetion
	db_connector.DbConnector()

	ctx := context.Background()

	// import the connector from db
	pool := config.GlobalDbConn.PgsqlConn

	// Step 6: Retrieve message data
	messages := pgsqlcrud.GetMessageData(ctx, "message", pool, "6LRcGlCjvNB", "Anny")

	// print the message data
	//fmt.Printf("Messages: %+v\n", messages)

	// to print the message from the database
	for _, m := range messages {
		fmt.Println("Message from %s to %s: %s\n",m.MessageID, m.Sender, m.Receiver, m.Message, m.Timestamp, m.Read)
	}

	
	

}