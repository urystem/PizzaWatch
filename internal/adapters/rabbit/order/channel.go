package order

import (
	"errors"

	"github.com/rabbitmq/amqp091-go"
)

func (r *rabbit) createChannel(dsn string) error {
	var err error
	r.rabbit, err = amqp091.Dial(dsn)
	if err != nil {
		return err
	}
	r.connClose = make(chan *amqp091.Error)
	r.rabbit.NotifyClose(r.connClose)
	r.orderCh, err = r.rabbit.Channel()
	if err != nil {
		return errors.Join(r.rabbit.Close(), err)
	}

	err = initKitchenQueue(r.orderCh)
	if err != nil {
		return errors.Join(r.rabbit.Close(), err)
	}
	r.notifyCh, err = r.rabbit.Channel()
	if err != nil {
		return errors.Join(r.rabbit.Close(), err)
	}

	err = initNotifyEx(r.notifyCh)
	if err != nil {
		return errors.Join(r.rabbit.Close(), err)
	}
	return nil
}
