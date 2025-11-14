package kitchen

import (
	"github.com/rabbitmq/amqp091-go"
)

const (
	orderExchange        = "orders_topic"
	notificationExchange = "notifications_fanout"
)

var queueNames = map[string]string{
	"dinein":   "kitchen_dine_in_queue",
	"takeout":  "kitchen_queue",
	"delivery": "kitchen_delivery",
}

func (r *rabbit) declareOrderTopic(orderCh *amqp091.Channel, types []string) error {
	err := orderCh.ExchangeDeclare(
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
	for _, v := range types {
		err := r.initQueue(orderCh, queueNames[v])
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *rabbit) initQueue(orderCh *amqp091.Channel, qName string) error {
	q, err := orderCh.QueueDeclare(qName, true, false, false, false, nil)
	if err != nil {
		return err
	}
	msgs, err := orderCh.Consume(
		q.Name,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	go r.merger(msgs)
	return nil
}

func (r *rabbit) merger(msgs <-chan amqp091.Delivery) {
	for msg := range msgs {
		r.ordersJobs <- NewQat(msg)
	}
}
