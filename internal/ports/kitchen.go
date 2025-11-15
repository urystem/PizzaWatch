package ports

import (
	"context"
	"pizza/internal/domain"
)

type KitchenPsql interface {
	CloseDB()
	CreateOrUpdateWorker(ctx context.Context, name string, types []string) ([]string, error)
	UpdateToOffline(ctx context.Context, name string) error
	AddOrderProcessed(ctx context.Context, name string) error
	UpdateStatusOrder(ctx context.Context, orderNumber, status, processedBy string) error
}

type KitchenRabbit interface {
	GiveChannel() <-chan QatJoldama
	CloseRabbit() error
	PublishNotify(ctx context.Context, zat *domain.LogMessageKitchen) error
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
