package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var DbConn *pgx.Conn

func ConnectDatabase() {
	dbConn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	DbConn = dbConn
	if err != nil {
		fmt.Printf("Unable to connect db: %v\n", err)
		os.Exit(1)
	}
}
