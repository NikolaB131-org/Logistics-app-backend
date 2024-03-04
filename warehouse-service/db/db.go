package db

import (
	"context"
	"fmt"
	"os"

	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/config"
	"github.com/jackc/pgx/v5"
)

var DbConn *pgx.Conn

func ConnectDatabase() {
	dbConn, err := pgx.Connect(context.Background(), config.Config.DbUrl)
	DbConn = dbConn
	if err != nil {
		fmt.Printf("Unable to connect db: %v\n", err)
		os.Exit(1)
	}
}
