package services

import (
	"context"
	"log/slog"
	"pizza/internal/domain"
	"pizza/internal/ports"
	"time"
)

type order struct {
	logg   *slog.Logger
	rabbit ports.OrderRabbit
	db     ports.OrderPsql
	sem    chan struct{} //concurent

}

func NewOrderService(logg *slog.Logger, rabbit ports.OrderRabbit, db ports.OrderPsql, maxConcurrent uint) ports.OrderUseCase {
	return &order{
		logg:   logg,
		db:     db,
		rabbit: rabbit,
		sem:    make(chan struct{}, maxConcurrent),
	}
}

func (u *order) CreateOrder(ctx context.Context, ord *domain.Order) (*domain.OrderStatus, error) {
	u.logg.Debug("receiving request", "action", "request_received")
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case u.sem <- struct{}{}:
		defer func() {
			<-u.sem
		}()
	}
	ordInsert := &domain.OrderPublish{Order: *ord}
	for _, item := range ord.Items {
		ordInsert.TotalAmount += item.Price * float64(item.Quantity)
	}

	if ordInsert.TotalAmount > 100 {
		ordInsert.Priority = 10
	} else if ordInsert.TotalAmount > 50 {
		ordInsert.Priority = 5
	} else {
		ordInsert.Priority = 1
	}

	if ord.OrderType != "dinein" {
		ord.TableNumber = nil
	}

	err := u.db.CreateOrder(ctx, ordInsert)
	if err != nil {
		return nil, err
	}

	err = u.rabbit.PublishOrder(ctx, ordInsert)
	if err != nil {
		return nil, err
	}

	err = u.rabbit.PublishNotify(ctx, &domain.OrderNotification{
		OrderNumber: ordInsert.OrderNumber,
		Status:      "received",
		TotalAmount: ordInsert.TotalAmount,
		Priority:    ordInsert.Priority,
		Timestamp:   time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return &domain.OrderStatus{
		OrderNumber: ordInsert.OrderNumber,
		Status:      "received",
		TotalAmount: ordInsert.TotalAmount,
	}, nil
}
