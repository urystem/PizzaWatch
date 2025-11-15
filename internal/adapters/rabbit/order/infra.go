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
	logger    *slog.Logger
	rabbit    *amqp091.Connection
	connClose chan *amqp091.Error
	// mu        sync.Mutex
	orderCh  *amqp091.Channel
	notifyCh *amqp091.Channel
	isClosed atomic.Bool
}

const (
	orderExchange        = "orders_topic"
	notificationExchange = "notifications_fanout"
)

func NewOrderRabbit(cfg config.CfgRabbitInter, slogger *slog.Logger) (ports.OrderRabbit, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.GetUser(), cfg.GetPassword(), cfg.GetHostName(), cfg.GetDBPort())
	myRab := &rabbit{logger: slogger}
	err := myRab.createChannel(dsn)
	if err != nil {
		return nil, err
	}
	go myRab.reconnectConn(dsn)
	return myRab, nil
}

func (r *rabbit) CloseRabbit() error {
	r.isClosed.Store(true)
	defer r.logger.Info("rabbit closed")
	return errors.Join(r.orderCh.Close(), r.notifyCh.Close(), r.rabbit.Close())
}

func (r *rabbit) PublishOrder(ctx context.Context, ord *domain.OrderPublish) error {
	b, err := json.Marshal(ord)
	if err != nil {
		return err
	}

	return r.orderCh.PublishWithContext(
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

func (r *rabbit) reconnectConn(url string) {
	for {
		<-r.connClose
		r.logger.Warn("rabbitMQ not working")
		for {
			if r.isClosed.Load() {
				return
			}
			r.logger.Info("trying to connect to rabbitmq")
			err := r.createChannel(url)
			if err != nil {
				time.Sleep(3 * time.Second)
				continue
			}
			r.logger.Info("connected to rabbitmq")
			break
		}
	}
}
