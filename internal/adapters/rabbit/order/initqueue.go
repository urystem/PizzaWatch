package order

import (
	"github.com/rabbitmq/amqp091-go"
)

func initKitchenQueue(ch *amqp091.Channel) error {
	err := ch.ExchangeDeclare(
		orderExchange, // имя exchange
		"topic",       // тип (direct, fanout, topic, headers)
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return err
	}
	q1, err := ch.QueueDeclare("kitchen_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.QueueBind(q1.Name, "kitchen.takeout.#", orderExchange, false, nil)
	if err != nil {
		return err
	}

	q2, err := ch.QueueDeclare("kitchen_dine_in_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.QueueBind(q2.Name, "kitchen.dinein.#", orderExchange, false, nil)
	if err != nil {
		return err
	}

	q3, err := ch.QueueDeclare("kitchen_delivery", true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.QueueBind(q3.Name, "kitchen.delivery.#", orderExchange, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func initNotifyEx(ch *amqp091.Channel) error {
	return ch.ExchangeDeclare(
		notificationExchange, // name
		"fanout",             // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
}
