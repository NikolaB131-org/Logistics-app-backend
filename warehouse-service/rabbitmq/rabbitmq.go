package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Channel           *amqp.Channel
	NotificationQueue amqp.Queue
}

func New(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect rabbitmq: %w", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	notificationQueue, err := channel.QueueDeclare(
		"notifications", // name
		false,           // durable
		false,           // autoDelete
		false,           // exclusive
		false,           // noWait
		nil,             // arguments
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		Channel:           channel,
		NotificationQueue: notificationQueue,
	}, nil
}

func (r *RabbitMQ) SendNotification(message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r.Channel.PublishWithContext(ctx,
		"",                       // exchange
		r.NotificationQueue.Name, // routing key
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}
