package rabbitmq

import (
	"context"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var channel *amqp.Channel
var notificationQueue amqp.Queue

func ConnectRabbitMQ() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
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
