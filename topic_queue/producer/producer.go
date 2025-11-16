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
	q, err := ch.QueueDeclare(
		"kitchen_queue",
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ch.QueueBind(q.Name, "" /*"#"*/, "topic_ex", false, nil)
	if err != nil {
		log.Fatal(err)
	}

	q1, err := ch.QueueDeclare(
		"kitchen_dine_in_queue",
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	err = ch.QueueBind(q1.Name, "user.#", "topic_ex", false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// err = ch.Publish(
	// 	"topic_ex",
	// 	"order.uвввp", // routing key
	// 	false,         // mandatory
	// 	false,         // immediate керек емес флаг
	// 	amqp091.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        []byte("order"),
	// 	},
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// егер ол жақта ешқандай queue болмаса хат автоматты түрде өшіп кетеді
	messages := map[string]string{
		"order.new":           "New order received",
		"order.cancel":        "Order cancelled",
		"user.signup":         "New user signed up",
		"user.update.profile": "User updated profile",
	}

	for key, msg := range messages {
		err = ch.Publish(
			"topic_ex",
			key,   // routing key
			false, // mandatory
			false, // immediate
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		// log.Printf("Sent %s -> %s\n", key, msg)
	}
}
