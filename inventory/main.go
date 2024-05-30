package main

import (
	"inventory/connections"
	"inventory/errors"
	"log"
)

func main() {
	products := map[string]int{
		"watch": 5,
		"shoes": 2,
	}

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
	errors.FailOnError(err, "Failed on declare a queue")

	err = ch.Qos(
		1,
		0,
		false,
	)
	errors.FailOnError(err, "Failed to set Qos")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false, // auto ack
		false,
		false,
		false,
		nil,
	)
	errors.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			element := string(d.Body)

			if products[element] > 0 {
				products[element]--
			} else {
				log.Println("this product is empty in inventory !")
				d.Ack(false)
				continue
			}
			log.Printf("Product %v amount: %v", element, products[element])

			d.Ack(false)
		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
