package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NikolaB131-org/logistics-backend/notifications-service/handlers"
	"github.com/NikolaB131-org/logistics-backend/notifications-service/rabbitmq"
)

func main() {
	rabbitmq.ConnectRabbitMQ()
	rabbitmq.ListenForNotifications()

	http.HandleFunc("/subscribe", handlers.HandleGetSubscribe)

	fmt.Printf("notifications-service is listening on port %s\n", os.Getenv("PORT"))
	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
	if err != nil {
		log.Fatal(err)
	}
}
