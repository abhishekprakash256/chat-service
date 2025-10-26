package config

import (
	"time"
)


type PgsqlDBConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
}

var PgsqlDefaultConfig = PgsqlDBConfig{
	Host:     "localhost",
	User:     "abhi",
	Password: "mysecretpassword",
	DBName:   "test_db",
	Port:     5432,
}




type LoginData struct {
	ChatID string
	UserOne string
	UserTwo string
}


type MessageOutData struct {
	messageid    int
	chatid       string
	sender   string
	receiver string
	message      string
	time    time.Time
	read         bool

}

type MessageData struct {
	MessageID    int
	ChatID       string
	Sender   string
	Receiver string
	Message      string
	Timestamp    time.Time
	Read         bool
}


// SQL to create the login table
var LoginTableSQL = `
CREATE TABLE IF NOT EXISTS login (
  chat_id     TEXT PRIMARY KEY,
  users_1     TEXT NOT NULL,
  users_2     TEXT NOT NULL
);`

// SQL to create the message table
var MessageTableSQL = `
CREATE TABLE IF NOT EXISTS message (
  message_id     SERIAL PRIMARY KEY,
  chat_id        TEXT NOT NULL REFERENCES login(chat_id) ON DELETE CASCADE,
  sender_name    TEXT NOT NULL,
  receiver_name  TEXT NOT NULL,
  message        TEXT NOT NULL,
  timestamp      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  read           BOOLEAN NOT NULL DEFAULT FALSE
);`



var MessageTable = "message"

var LoginTable = "login"