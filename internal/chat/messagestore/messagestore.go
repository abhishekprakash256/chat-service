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


// SaveMessage persists a chat message into the PostgreSQL database.
//
// It creates a MessageData struct with sender, receiver, message content,
// chat ID, current timestamp, and read status, then attempts to insert
// it into the "message" table via pgsqlcrud.InsertMessageData.
//
// Parameters:
//   - Sender   (string): Username or ID of the message sender.
//   - Receiver (string): Username or ID of the message recipient.
//   - Message  (string): The actual message text.
//   - ChatID   (string): Unique identifier of the chat session.
//
// Returns:
//   - error: Returns an error if inserting into the database fails,
//            otherwise returns nil on success.
//
// Behavior:
//   - On success, logs "Message saved successfully in SaveMessage".
//   - On failure, logs the error and returns a descriptive error value.
//
// Example:
//   err := messagestore.SaveMessage("alice", "bob", "Hello Bob!", "chat123")
//   if err != nil {
//       log.Println("Failed to save message:", err)
//   }
//
func SaveMessage(Sender string, Receiver string, Message string, ChatID string) (int64, time.Time, error) {
	pool := config.GlobalDbConn.PgsqlConn
	ctx := context.Background()

	msg := config.MessageData{
		ChatID:    ChatID,
		Sender:    Sender,
		Receiver:  Receiver,
		Message:   Message,
		Timestamp: time.Now(),
		Read:      false,
	}

	messageID, msgTime ,err := pgsqlcrud.InsertMessageData(ctx, config.MessageTable, pool, msg)
	if err != nil {
		log.Println("failing")
		return 0, time.Time{}, fmt.Errorf("failed to insert message into DB: %v", err)
	}

	log.Printf("Message saved successfully with ID: %d", messageID)
	return messageID, msgTime , nil
}

