package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NikolaB131-org/logistics-backend/notifications-service/handlers"
	"github.com/NikolaB131-org/logistics-backend/notifications-service/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var channel *amqp.Channel
var notificationQueue amqp.Queue

func ConnectRabbitMQ() {
	conn, err := amqp.Dial(config.Config.RabbitmqUrl)
	if err != nil {
		fmt.Printf("Unable to connect rabbitmq: %v\n", err)
		os.Exit(1)
	}
	channel, _ = conn.Channel()
	notificationQueue, _ = channel.QueueDeclare(
		"notifications", // name
		false,           // durable
		false,           // autoDelete
		false,           // exclusive
		false,           // noWait
		nil,             // arguments
	)
}

func SendNotification(message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	channel.PublishWithContext(ctx,
		"",                     // exchange
		notificationQueue.Name, // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func ListenForNotifications() {
	msgs, _ := channel.Consume(
		notificationQueue.Name, // queue
		"",                     // consumer
		true,                   // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			for _, client := range handlers.Clients {
				client <- string(d.Body)
			}
		}
	}()
}
