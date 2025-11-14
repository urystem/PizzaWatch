package kitchen

import (
	"errors"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func (r *rabbit) createChannel(dsn string, orderTypes []string, prefetch int) error {
	var err error
	r.conn, err = amqp091.Dial(dsn)
	if err != nil {
		return err
	}
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
	// r.notifyCh, err = r.rabbit.Channel()
	// if err != nil {
	// 	return errors.Join(r.rabbit.Close(), err)
	// }

	// err = initNotifyEx(r.notifyCh)
	// if err != nil {
	// 	return errors.Join(r.rabbit.Close(), err)
	// }
	return nil
}

func (r *rabbit) reconnectConn(url string, ords []string, prefetch int) {
	for {
		<-r.connClose
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
