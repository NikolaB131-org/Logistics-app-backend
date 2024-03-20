package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NikolaB131-org/logistics-backend/warehouse-service/db"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/handlers"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/rabbitmq"
)

func main() {
	db.ConnectDatabase()
	rabbitmq.ConnectRabbitMQ()

	http.HandleFunc("/products", handlers.HandleProducts)

	fmt.Printf("warehouse-service is listening on port %s\n", os.Getenv("PORT"))
	err := http.ListenAndServe(os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
