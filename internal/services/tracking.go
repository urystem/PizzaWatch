package services

import (
	"context"
	"log/slog"

	"pizza/internal/domain"
	"pizza/internal/ports"
)

type tracing struct {
	logg *slog.Logger
	db   ports.TrackingSQL
}

func NewTrackingService(logg *slog.Logger, db ports.TrackingSQL) ports.TrackingUse {
	return &tracing{
		logg: logg,
		db:   db,
	}
}

func (t *tracing) GetWorkersStatus(ctx context.Context, heartbeatInterval uint) ([]domain.WorkerStatus, error) {
	return t.db.GetWorkers(ctx, heartbeatInterval)
}

func (t *tracing) OrderStatusUpdate(ctx context.Context, number string) (*domain.OrderStatusUpdate, error) {
	return t.db.OrderStatusUpdate(ctx, number)
}

func (t *tracing) GetOrderHistory(ctx context.Context, number string) ([]domain.OrderStatusEvent, error) {
	return t.db.GetOrderHistory(ctx, number)
}
