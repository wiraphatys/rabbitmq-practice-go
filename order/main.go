package main

import (
	"context"
	"log"
	"order/connections"
	"order/errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := connections.NewRabbitMQConnection()
	if err != nil {
		log.Fatalf("failed to create a RabbitMQ connection: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	errors.FailOnError(err, "Failed on declare a queue")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"order_created",
		true,
		false,
		false,
		false,
		nil,
	)
	errors.FailOnError(err, "Failed to declare a queue")

	var val []string = []string{"shoes", "shoes", "shoes"}

	for _, e := range val {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body := e
		err = ch.PublishWithContext(ctx,
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			},
		)
		errors.FailOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", body)
	}
}
