package kitchen

import (
	"fmt"
	"log/slog"
	"pizza/internal/config"
	"pizza/internal/ports"
	"sync/atomic"

	"github.com/rabbitmq/amqp091-go"
)

type rabbit struct {
	logger *slog.Logger
	conn   *amqp091.Connection
	// Orderch    *amqp091.Channel
	connClose  chan *amqp091.Error
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
	return r.conn.Close()
}
