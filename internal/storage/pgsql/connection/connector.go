// to make the connetor of the pgsql 

package pgsql

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)


func ConnectPgSql(host, userName, password, dbName string, port int) (*pgxpool.Pool, error) {
	// Format the connection string properly

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		userName, password, host, port, dbName)

	// Connect to PostgreSQL
	pool, err := pgxpool.New(context.Background(), connStr)
	
	// get the error
	if err != nil {
		log.Fatal(" Failed to connect:", err)
	}

	//defer conn.Close(context.Background())
	fmt.Println("Connected to PostgreSQL with pgx!")

	return pool, nil
}