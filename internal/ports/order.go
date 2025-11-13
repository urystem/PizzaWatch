package ports

import (
	"context"
	"pizza/internal/domain"
)

type OrderUseCase interface {
	CreateOrder(ctx context.Context, ord *domain.Order) (*domain.OrderStatus, error)
}

type OrderPsql interface {
	CloseDB()
	CreateOrder(ctx context.Context, ord *domain.OrderPublish) error
}

type OrderRabbit interface {
	CloseRabbit() error
	PublishOrder(ctx context.Context, ord *domain.OrderPublish) error
	PublishNotify(ctx context.Context, ord *domain.OrderNotification) error
}

type ServerInter interface {
	StartServer() error
	ShutDownServer(ctx context.Context) error
}
