package kitchen

import (
	"encoding/json"
	"fmt"
	"pizza/internal/domain"
	"pizza/internal/ports"

	"github.com/rabbitmq/amqp091-go"
)

type qat struct {
	hat amqp091.Delivery
}

func NewQat(hat amqp091.Delivery) ports.QatJoldama {
	return &qat{
		hat: hat,
	}
}

func (q *qat) GiveBody() (*domain.OrderPublish, error) {
	ord := new(domain.OrderPublish)
	err := json.Unmarshal(q.hat.Body, ord)
	if err != nil {
		return nil, err
	}
	return ord, nil
}

func (q *qat) Qaitar() error {
	return q.hat.Nack(false, true)
}

func (q *qat) Joi() error {
	return q.hat.Nack(false, false)
}

func (q *qat) Rastau() error {
	defer fmt.Println("rautaldy")
	return q.hat.Ack(false)
}
