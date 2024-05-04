package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var DbConn *pgx.Conn

func ConnectDatabase() error {
	dbConn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	DbConn = dbConn
	if err != nil {
		return fmt.Errorf("unable to connect db: %w", err)
	}
	return nil
}
