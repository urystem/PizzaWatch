package main

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp091.Dial("amqp://urystem:admin123@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ch.ExchangeDeclare(
		"topic_example", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	q, err := ch.QueueDeclare(
		"ddd", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	err = ch.QueueBind(
		q.Name,
		"black",
		"topic_example",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	for d := range msgs {
		log.Printf(" [x] %s", d.Body)
	}
}
