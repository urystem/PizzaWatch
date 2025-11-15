package kitchen

import (
	"errors"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func (r *rabbit) createChannel(dsn string, orderTypes []string, prefetch int) error {
	myConn, err := amqp091.Dial(dsn)
	if err != nil {
		return err
	}
	r.conn = myConn
	r.connClose = make(chan *amqp091.Error)
	r.conn.NotifyClose(r.connClose)
	orderCh, err := r.conn.Channel()
	if err != nil {
		return errors.Join(r.conn.Close(), err)
	}
	err = orderCh.Qos(prefetch, 0, false)
	if err != nil {
		return err
	}
	err = r.declareOrderTopic(orderCh, orderTypes)
	if err != nil {
		return errors.Join(r.conn.Close(), err)
	}

	r.notifyCh, err = r.conn.Channel()
	if err != nil {
		return errors.Join(r.conn.Close(), err)
	}
	err = r.notifyCh.ExchangeDeclare(
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
	// if err != nil {
	// 	return errors.Join(r.rabbit.Close(), err)
	// }
	return nil
}

func (r *rabbit) reconnectConn(url string, ords []string, prefetch int) {
	for {
		<-r.connClose
		if r.isClosed.Load() {
			return
		}
		r.logger.Warn("rabbitMQ not working")
		for {
			if r.isClosed.Load() {
				return
			}
			r.logger.Info("trying to connect to rabbitmq")
			err := r.createChannel(url, ords, prefetch)
			if err != nil {
				time.Sleep(3 * time.Second)
				continue
			}
			r.logger.Info("connected to rabbitmq")
			break
		}
	}
}
