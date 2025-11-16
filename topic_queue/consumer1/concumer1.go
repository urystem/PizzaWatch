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
	defer ch.Close()
	err = ch.ExchangeDeclare(
		"topic_ex", // name
		"topic",    // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Создаём очередь
	q, err := ch.QueueDeclare(
		"kitchen_dine_in_queue", // имя очереди
		true,                    // durable
		false,                   // autoDelete
		false,                   // exclusive
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ch.QueueBind(q.Name, "#", "topic_ex", false, nil)
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	for d := range msgs {
		log.Printf("Received: [%s] %s\n", d.RoutingKey, d.Body)
	}
}
