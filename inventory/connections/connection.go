package connections

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConnection() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://root:1234@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	return conn, nil
}
