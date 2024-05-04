package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var DbConn *pgx.Conn

func ConnectDatabase(url string) error {
	dbConn, err := pgx.Connect(context.Background(), url)
	DbConn = dbConn
	if err != nil {
		return fmt.Errorf("unable to connect db: %w", err)
	}
	return nil
}
