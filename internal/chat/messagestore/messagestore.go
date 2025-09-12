/*
the package to save the message from chat 
from user using pgsql 

*/

package messagestore

import (

	"time"
	"context"
	"fmt"
	"log"

	"chat-service/internal/config"
	pgsqlcrud "chat-service/internal/storage/pgsql/crud"

	
)



func SaveMessage(Sender string , Receiver string , Message string , ChatID string) error {

	// get the pgsql conn 
	pool := config.GlobalDbConn.PgsqlConn

	// start the context 
	ctx := context.Background()


	// Test message data
	msg := config.MessageData{
		ChatID:   ChatID,
		Sender:   Sender,
		Receiver: Receiver,
		Message:  Message,
		Timestamp: time.Now(),
		Read:  false,
	}

	// Insert the message data
	if !pgsqlcrud.InsertMessageData(ctx, "message", pool, msg) {

		return fmt.Errorf("failed to insert message into DB")
	}

	log.Println("Message saved succesfully in savemessage")

	return nil


}