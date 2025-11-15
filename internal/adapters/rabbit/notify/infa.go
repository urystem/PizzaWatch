package notify

import (
	"fmt"
	"log/slog"
	"pizza/internal/config"
	"pizza/internal/ports"
	"sync/atomic"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type rabbit struct {
	logger    *slog.Logger
	conn      *amqp091.Connection
	connClose chan *amqp091.Error
	notiJobs  chan []byte
	isClosed  atomic.Bool
}

func NewNotifyRabbit(cfg config.CfgRabbitInter, logg *slog.Logger) (ports.NotifyRabbit, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.GetUser(), cfg.GetPassword(), cfg.GetHostName(), cfg.GetDBPort())
	myRab := &rabbit{
		logger:   logg,
		notiJobs: make(chan []byte),
	}
	err := myRab.createChannel(dsn)
	if err != nil {
		return nil, err
	}
	go myRab.reconnectConn(dsn)
	return myRab, nil
}

func (r *rabbit) reconnectConn(url string) {
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

func (r *rabbit) GiveChannel() <-chan []byte {
	return r.notiJobs
}

func (r *rabbit) CloseRabbit() error {
	r.isClosed.Store(true)
	defer r.logger.Info("rabbit closed")
	return r.conn.Close()
}
