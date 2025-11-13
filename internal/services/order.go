package services

import (
	"context"
	"pizza/internal/domain"
	"pizza/internal/ports"
	"time"
)

type order struct {
	rabbit ports.OrderRabbit
	db     ports.OrderPsql
}

func NewOrderService(rabbit ports.OrderRabbit, db ports.OrderPsql) ports.OrderUseCase {
	return &order{
		db:     db,
		rabbit: rabbit,
	}
}

func (u *order) CreateOrder(ctx context.Context, ord *domain.Order) (*domain.OrderStatus, error) {
	ordInsert := &domain.OrderPublish{Order: *ord}
	for _, item := range ord.Items {
		ordInsert.TotalAmount += item.Price
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
