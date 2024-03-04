package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NikolaB131-org/logistics-backend/notifications-service/handlers"
	"github.com/NikolaB131-org/logistics-backend/notifications-service/internal/config"
	"github.com/NikolaB131-org/logistics-backend/notifications-service/rabbitmq"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	rabbitmq.ConnectRabbitMQ()
	rabbitmq.ListenForNotifications()

	http.HandleFunc("/subscribe", handlers.HandleGetSubscribe)

	fmt.Printf("notifications-service is listening port %d\n", config.Config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
