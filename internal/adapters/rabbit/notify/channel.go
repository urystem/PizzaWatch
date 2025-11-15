package notify

import (
	"github.com/rabbitmq/amqp091-go"
)

func (r *rabbit) createChannel(url string) error {
	var err error
	r.conn, err = amqp091.Dial(url)
	if err != nil {
		return err
	}
	r.connClose = make(chan *amqp091.Error)
	r.conn.NotifyClose(r.connClose)
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	err = ch.ExchangeDeclare(
		"notifications_fanout", // name
		"fanout",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return err
	}
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,                 // queue name
		"",                     // routing key
		"notifications_fanout", // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}
	go r.merger(msgs)
	return nil
}

func (r *rabbit) merger(msgs <-chan amqp091.Delivery) {
	for m := range msgs {
		go func() {
			r.notiJobs <- m.Body
		}()
	}
}
