// the function to make the table and schema 


package pgsql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateSchema takes custom SQL strings for login and message table creation,
// and uses the provided context properly.
func CreateSchema(ctx context.Context, pgconnector *pgxpool.Pool, loginSQL, messageSQL string) error {
	if _, err := pgconnector.Exec(ctx, loginSQL); err != nil {
		return fmt.Errorf("failed to create login table: %w", err)
	}

	if _, err := pgconnector.Exec(ctx, messageSQL); err != nil {
		return fmt.Errorf("failed to create message table: %w", err)
	}

	fmt.Println("Tables created (or already exist).")
	return nil
}