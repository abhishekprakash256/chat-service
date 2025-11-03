/*
The message reader is the file to read the message 

*/

package messagereader


import (
	"log"
	"fmt"
	//"time"
	"encoding/json"
	"github.com/gorilla/websocket"
	"chat-service/internal/config"
	"chat-service/internal/chat/messagestore"
)


// ReadMessage continuously listens for new WebSocket messages from a client.
//
// Workflow:
//   1. Reads raw messages from the WebSocket connection.
//   2. Pushes the raw message bytes into the global broadcast channel.
//   3. Logs the raw message for debugging.
//   4. Attempts to parse the message into config.IncomingMessage (JSON format).
//   5. Logs routing information (sender, receiver, chat ID).
//   6. Persists the message into PostgreSQL using messagestore.SaveMessage.
//
// Parameters:
//   - conn (*websocket.Conn): Active WebSocket connection for the client.
//
// Behavior:
//   - If the WebSocket connection closes or an error occurs while reading,
//     the function logs the error and exits.
//   - If a message cannot be unmarshaled into IncomingMessage, the error is logged
//     and the connection is closed.
//   - If saving the message to PostgreSQL fails, the error is logged
//     and the connection is closed.
//   - On successful save, a confirmation log entry is written.
//
// Note:
//   - Currently, any error (read, parse, or save) will break out of the loop and
//     terminate the connection. Consider replacing `return` with `continue` in some
//     error branches if you want the connection to persist despite bad messages.
//
// Example:
//   go ReadMessage(conn) // Run as a goroutine for each client connection
//
func ReadMessage(conn *websocket.Conn) {
	for {
		// Step 1: Read message from WebSocket
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			return
		}

		// ping pong for client
		
		var raw map[string]interface{}
		if err := json.Unmarshal(msg, &raw); err == nil {
			// Handle ping from client
			if raw["type"] == "ping" {

				log.Println("Ping from Client" ,raw["sessionID"])
				_ = conn.WriteJSON(map[string]string{"type": "pong"})
				continue
			}
		}
		
		// Step 2: Broadcast raw message
		//config.BroadCast <- msg

		// Step 3: Debug log
		fmt.Println("Raw message:", string(msg))

		// Step 4: Parse into IncomingMessage
		var incoming config.IncomingMessage
		if err := json.Unmarshal(msg, &incoming); err != nil {
			log.Println("Invalid message format:", err)
			return
		}

		// Step 6: Save message to PostgreSQL
		messageID, msgTime , err := messagestore.SaveMessage(
			incoming.Sender,
			incoming.Receiver,
			incoming.Message,
			incoming.ChatID,
		)

		if err != nil {
			log.Printf("Failed to save message for chat %s: %v", incoming.ChatID, err)
			continue
		}

		log.Println("Message saved successfully in ReadMessage")
		
		// prepare out going message
		// Step 4: Prepare outgoing message
		outgoingmessage := config.OutgoingMessage{
			MessageID: messageID,
			ChatID:    incoming.ChatID,
			Sender:    incoming.Sender,
			Receiver:  incoming.Receiver,
			Message:   incoming.Message,
			Timestamp:  msgTime ,  // add the time stamp from the db
		}

		outgoingBytes, err := json.Marshal(outgoingmessage)
		if err != nil {
			log.Println("Error marshaling outgoing message:", err)
			continue
		}

		// Step 5: Send to broadcast channel
		config.BroadCast <- outgoingBytes


		
	}
}

