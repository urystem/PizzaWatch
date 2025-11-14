package ports

import (
	"context"
	"pizza/internal/domain"
)

type KitchenPsql interface {
	CloseDB()
	CreateOrUpdateWorker(ctx context.Context, name string, types []string) error
	AddOrderProcessed(ctx context.Context, name string) error
}

type KitchenRabbit interface {
	GiveChannel() <-chan QatJoldama
	CloseRabbit() error
}

type QatJoldama interface {
	GiveBody() (*domain.OrderPublish, error)
	Qaitar() error
	Rastau() error
	Joi() error
}

type KitchenService interface {
	StartWork()
}
