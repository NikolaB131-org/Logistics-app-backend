package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NikolaB131-org/logistics-backend/warehouse-service/db"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/handlers"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/config"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/rabbitmq"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db.ConnectDatabase()
	rabbitmq.ConnectRabbitMQ()

	http.HandleFunc("/products", handlers.HandleProducts)

	fmt.Printf("warehouse-service is listening port %d\n", config.Config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
