package ports

import (
	"context"

	"pizza/internal/domain"
)

type TrackingSQL interface {
	CloseDB()
	GetWorkers(ctx context.Context, heartbeatInterval uint) ([]domain.WorkerStatus, error)
	OrderStatusUpdate(ctx context.Context, number string) (*domain.OrderStatusUpdate, error)
	GetOrderHistory(ctx context.Context, number string) ([]domain.OrderStatusEvent, error)
}

type TrackingUse interface {
	GetWorkersStatus(ctx context.Context, heartbeatInterval uint) ([]domain.WorkerStatus, error)
	OrderStatusUpdate(ctx context.Context, number string) (*domain.OrderStatusUpdate, error)
	GetOrderHistory(ctx context.Context, number string) ([]domain.OrderStatusEvent, error)
}
