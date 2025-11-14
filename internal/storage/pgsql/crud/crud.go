// to make the crud operation on the data 

/*
sample data 
| message_id | chat_id | sender_name | receiver_name | message | timestamp          | read |
|------------|---------|-------------|----------------|---------|---------------------|------|
| ...        | abc123  | "Abhi"      | "Anny"         | "Hello" | 2025-07-06 15:00:00 | TRUE |

*/


package pgsql

import (
	"context"
	"fmt"
	"time"
	"chat-service/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)




func GetMessageData(ctx context.Context, tableName string, pgconnector *pgxpool.Pool, chatID string, user string) []config.MessageData {

	query := fmt.Sprintf(`
		SELECT * FROM (
			SELECT message_id, chat_id, sender_name, receiver_name, message, timestamp, read
			FROM %s
			WHERE chat_id = $1 AND (sender_name = $2 OR receiver_name = $2)
			ORDER BY timestamp DESC
			LIMIT 100
		) AS sub
		ORDER BY timestamp ASC
	`, tableName)

	rows, err := pgconnector.Query(ctx, query, chatID, user)
	if err != nil {
		fmt.Println("Query failed:", err)
		return nil
	}
	defer rows.Close()

	var messages []config.MessageData
	for rows.Next() {
		var msg config.MessageData
		// chnages here 
		err := rows.Scan(&msg.MessageID, &msg.ChatID, &msg.Sender, &msg.Receiver, &msg.Message, &msg.Timestamp, &msg.Read)
		if err != nil {
			fmt.Println("Row scan failed:", err)
			continue
		}
		messages = append(messages, msg)
	}

	return messages
}


// the function to fetch the message data using the messageid
// the function takes the json as chatId , username and messageid
// gets the message from the messageid and message before
// needs to be changed as per frontend reqs
func GetMessageDataID(ctx context.Context, tableName string, pgconnector *pgxpool.Pool, chatID string, user string , messageid int) []config.MessageData {

	query := fmt.Sprintf(`
		SELECT * FROM (
			SELECT message_id, chat_id, sender_name, receiver_name, message, timestamp, read
			FROM %s
			WHERE chat_id = $1 
			AND (sender_name = $2 OR receiver_name = $2)
			AND message_id < $3
			ORDER BY timestamp DESC
			LIMIT 50
		) AS sub
		ORDER BY timestamp ASC
	`, tableName)

	rows, err := pgconnector.Query(ctx, query, chatID, user , messageid)
	if err != nil {
		fmt.Println("Query failed:", err)
		return nil
	}
	defer rows.Close()

	var messages []config.MessageData
	for rows.Next() {
		var msg config.MessageData
		// chnages here 
		err := rows.Scan(&msg.MessageID, &msg.ChatID, &msg.Sender, &msg.Receiver, &msg.Message, &msg.Timestamp, &msg.Read)
		if err != nil {
			fmt.Println("Row scan failed:", err)
			continue
		}
		messages = append(messages, msg)
	}

	return messages
}


func GetLoginData(ctx context.Context, tableName string, pgconnector *pgxpool.Pool, chatID string) (config.LoginData, error) {
	query := fmt.Sprintf(`SELECT chat_id, users_1, users_2 FROM %s WHERE chat_id = $1`, tableName)

	var data config.LoginData
	err := pgconnector.QueryRow(ctx, query, chatID).Scan(&data.ChatID, &data.UserOne, &data.UserTwo)
	if err != nil {
		return config.LoginData{}, fmt.Errorf("login data not found: %w", err)
	}

	return data, nil
}



// InsertLoginData inserts a row into the login table.
func InsertLoginData(ctx context.Context, tableName string, pgconnector *pgxpool.Pool, data config.LoginData) bool {
	query := fmt.Sprintf(`
		INSERT INTO %s (chat_id, users_1, users_2)
		VALUES ($1, $2, $3)
		ON CONFLICT (chat_id) DO NOTHING
	`, tableName)

	_, err := pgconnector.Exec(ctx, query, data.ChatID, data.UserOne, data.UserTwo)
	

	if err != nil {
		fmt.Println("Insert into login failed:", err)
		return false
	}

	fmt.Println("Login inserted (or already exists).")
	return true
}


/*
func InsertMessageData(ctx context.Context, tableName string, pgconnector *pgxpool.Pool, data config.MessageData) bool {

	// insert the data into the message table

	query := fmt.Sprintf(`
		INSERT INTO %s (chat_id, sender_name, receiver_name, message, timestamp, read)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, tableName)

	_, err := pgconnector.Exec(
		ctx,
		query,
		data.ChatID,
		data.Sender,
		data.Receiver,
		data.Message,
		data.Timestamp,
		data.Read,
	)

	if err != nil {
		fmt.Println("Insert failed:", err)
		return false
	}

	fmt.Println("Message inserted")
	return true
}
*/

func InsertMessageData(ctx context.Context, tableName string, pgconnector *pgxpool.Pool, data config.MessageData) (int64, time.Time , error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (chat_id, sender_name, receiver_name, message, timestamp, read)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING message_id , timestamp
	`, tableName)

	var messageID int64
	var msgTime time.Time

	err := pgconnector.QueryRow(
		ctx,
		query,
		data.ChatID,
		data.Sender,
		data.Receiver,
		data.Message,
		data.Timestamp,
		data.Read,
	).Scan(&messageID, &msgTime)

	if err != nil {
		fmt.Printf("Insert failed: %v\n", err)
		return 0, time.Time{} , err
	}

	fmt.Printf("Message inserted with ID: %d\n", messageID)
	return messageID, msgTime, nil
}



func DeleteLoginData(ctx context.Context, tableName string, pgconnector *pgxpool.Pool, chatID string) bool {

	// delete the message per id

	query := fmt.Sprintf(`DELETE FROM %s WHERE chat_id = $1`, tableName)

	_, err := pgconnector.Exec(ctx, query, chatID)
	
	if err != nil {
		fmt.Println("Delete failed:", err)
		return false
	}

	fmt.Println("Login data deleted for chat_id:", chatID)
	return true
}


func DeleteMessageData(ctx context.Context, tableName string, pgconnector *pgxpool.Pool, chatID string) bool {

	// delete the message per id

	query := fmt.Sprintf(`DELETE FROM %s WHERE chat_id = $1`, tableName)

	_, err := pgconnector.Exec(ctx, query, chatID)

	if err != nil {
		fmt.Println("Delete failed:", err)
		return false
	}

	fmt.Println("Messages deleted for chat_id:", chatID)
	
	return true
}