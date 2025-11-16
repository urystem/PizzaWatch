package kitchen

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync/atomic"

	"pizza/internal/config"
	"pizza/internal/domain"
	"pizza/internal/ports"

	"github.com/rabbitmq/amqp091-go"
)

type rabbit struct {
	logger *slog.Logger
	conn   *amqp091.Connection
	// Orderch    *amqp091.Channel
	connClose  chan *amqp091.Error
	notifyCh   *amqp091.Channel
	ordersJobs chan ports.QatJoldama
	isClosed   atomic.Bool
}

func NewKitchenRabbit(cfg config.CfgRabbitInter, slogger *slog.Logger, prefetch int, types []string) (ports.KitchenRabbit, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.GetUser(), cfg.GetPassword(), cfg.GetHostName(), cfg.GetDBPort())
	myRab := &rabbit{
		logger:     slogger,
		ordersJobs: make(chan ports.QatJoldama),
	}

	err := myRab.createChannel(dsn, types, prefetch)
	if err != nil {
		return nil, err
	}

	go myRab.reconnectConn(dsn, types, prefetch)
	return myRab, nil
}

func (r *rabbit) GiveChannel() <-chan ports.QatJoldama {
	return r.ordersJobs
}

func (r *rabbit) CloseRabbit() error {
	r.isClosed.Store(true)
	defer r.logger.Info("rabbit closed")
	return r.conn.Close()
}

func (r *rabbit) PublishNotify(ctx context.Context, zat *domain.LogMessageKitchen) error {
	b, err := json.Marshal(zat)
	if err != nil {
		return err
	}
	return r.notifyCh.PublishWithContext(
		ctx,
		"notifications_fanout",
		"",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	)
}
