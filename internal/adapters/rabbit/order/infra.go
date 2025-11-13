package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"pizza/internal/config"
	"pizza/internal/domain"
	"pizza/internal/ports"
	"sync/atomic"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type rabbit struct {
	rabbit    *amqp091.Connection
	connClose chan *amqp091.Error
	orderCh   *amqp091.Channel
	notifyCh  *amqp091.Channel
	sem       chan struct{} //concurent
	isClosed atomic.Bool
}

const (
	orderExchange        = "orders_topic"
	notificationExchange = "notifications_fanout"
)

func (r *rabbit) reconnectConn(url string) {
	for {
		<-r.connClose
		for {
			if r.isClosed.Load() {
				return
			}
			slog.Error("try")
			conn, err := amqp091.Dial(url)
			if err != nil {
				time.Sleep(3 * time.Second)
				continue
			}
			// r.connClose = make(chan *amqp091.Error)
			conn.NotifyClose(r.connClose)
			r.rabbit = conn
			break
		}
	}
}

func NewOrderRabbit(cfg config.CfgRabbitInter, maxConcurrent uint) (ports.OrderRabbit, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.GetUser(), cfg.GetPassword(), cfg.GetHostName(), cfg.GetDBPort())
	conn, err := amqp091.Dial(dsn)
	if err != nil {
		return nil, err
	}
	ch1, err := conn.Channel()
	if err != nil {
		return nil, errors.Join(conn.Close(), err)
	}
	err = initKitchenQueue(ch1)
	if err != nil {
		return nil, errors.Join(ch1.Close(), conn.Close(), err)
	}
	ch2, err := conn.Channel()
	if err != nil {
		return nil, errors.Join(ch1.Close(), conn.Close(), err)
	}

	err = initNotifyEx(ch2)
	if err != nil {
		return nil, errors.Join(ch1.Close(), ch2.Close(), conn.Close(), err)
	}
	myRab := &rabbit{
		rabbit:   conn,
		orderCh:  ch1,
		notifyCh: ch2,
		sem:      make(chan struct{}, maxConcurrent),
	}
	go myRab.reconnectConn(dsn)
	return myRab, nil
}

func (r *rabbit) CloseRabbit() error {
	r.isClosed.Store(true)
	return errors.Join(r.orderCh.Close(), r.notifyCh.Close(), r.rabbit.Close())
}

func (r *rabbit) PublishOrder(ctx context.Context, ord *domain.OrderPublish) error {
	b, err := json.Marshal(ord)
	if err != nil {
		return err
	}
	select {
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout")
	case <-ctx.Done():
		return ctx.Err()
	case r.sem <- struct{}{}:
		err := r.orderCh.PublishWithContext(
			ctx,
			orderExchange,
			fmt.Sprintf("kitchen.%s.%d", ord.OrderType, ord.Priority),
			false,
			false,
			amqp091.Publishing{
				ContentType:  "application/json",
				Body:         b,
				Priority:     ord.Priority,
				DeliveryMode: amqp091.Persistent, // 2 → persistent
			},
		)
		<-r.sem
		return err
	}
}

func (r *rabbit) PublishNotify(ctx context.Context, ord *domain.OrderNotification) error {
	b, err := json.Marshal(ord)
	if err != nil {
		return err
	}
	return r.notifyCh.PublishWithContext(
		ctx,
		notificationExchange,
		"",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	)
}
