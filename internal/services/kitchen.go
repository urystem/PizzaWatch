package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"pizza/internal/domain"
	"pizza/internal/ports"
)

type kitchen struct {
	ctx      context.Context
	slogger  *slog.Logger
	rabbit   ports.KitchenRabbit
	db       ports.KitchenPsql
	name     string
	hostName string
	types    map[string]struct{}
}

func NewKitchenService(ctx context.Context, slogger *slog.Logger, rabbit ports.KitchenRabbit, db ports.KitchenPsql, name, hostname string, types []string) (ports.KitchenService, error) {
	originalTypes, err := db.CreateOrUpdateWorker(context.Background(), name, types)
	if err != nil {
		return nil, err
	}
	typesMap := make(map[string]struct{})
	for _, tip := range originalTypes {
		_, ok := typesMap[tip]
		if ok {
			return nil, fmt.Errorf("%s", "dulpicated order type in sql")
		}
		typesMap[tip] = struct{}{}
	}
	return &kitchen{
		ctx:      ctx,
		slogger:  slogger,
		rabbit:   rabbit,
		db:       db,
		name:     name,
		types:    typesMap,
		hostName: hostname,
	}, nil
}

func (u *kitchen) StartWork() {
	jobs := u.rabbit.GiveChannel()
	defer u.slogger.Info("service stopped")
	for {
		select {
		case <-u.ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			go u.worker(job)
		}
	}
}

func (u *kitchen) worker(job ports.QatJoldama) error {
	order, err := job.GiveBody()
	if err != nil {
		return job.Joi()
	}
	loggNoti := &domain.LogMessageKitchen{
		Service:  "kitchen",
		Level:    "INFO",
		Hostname: u.hostName,
		Details: domain.Details{
			OrderNumber: order.OrderNumber,
		},
	}
	_, ok := u.types[order.OrderType]
	if !ok {
		loggNoti.Timestamp = time.Now()
		loggNoti.Level = "DEBUG"
		loggNoti.Action = fmt.Sprintf("nack the order : %s", order.OrderNumber)
		loggNoti.Message = "order type is not specialized"
		defer u.rabbit.PublishNotify(u.ctx, loggNoti)
		return job.Qaitar()
	}

	u.slogger.Debug("order processing", "action", "order_processing_started", "order_number", order.OrderNumber)
	err = u.db.UpdateStatusOrder(u.ctx, order.OrderNumber, "cooking", u.name)
	if err != nil {
		u.slogger.Error("cannot update order status", "error", err)
		loggNoti.Timestamp = time.Now()
		loggNoti.Level = "ERROR"
		loggNoti.Action = fmt.Sprintf("nack the order : %s", order.OrderNumber)
		loggNoti.Message = "cannot update status in sql"
		defer u.rabbit.PublishNotify(u.ctx, loggNoti)
		return job.Qaitar()
	}
	var duration time.Duration
	switch order.OrderType {
	case "dinein":
		duration = 8 * time.Second
	case "takeout":
		duration = 10 * time.Second
	case "delivery":
		duration = 12 * time.Second
	}
	select {
	case <-u.ctx.Done():
		defer u.db.UpdateStatusOrder(u.ctx, order.OrderNumber, "received", u.name)
		loggNoti.Timestamp = time.Now()
		loggNoti.Level = "INFO"
		loggNoti.Action = fmt.Sprintf("nack the order : %s", order.OrderNumber)
		loggNoti.Message = "Grasful shutdown"
		defer u.rabbit.PublishNotify(u.ctx, loggNoti)
		return errors.Join(err, job.Qaitar())
	case <-time.After(duration):
		err = u.db.UpdateStatusOrder(u.ctx, order.OrderNumber, "ready", u.name)
		if err != nil {
			u.slogger.Error("cannot update status of order", "error", err)
			return job.Qaitar()
		}
		err = u.db.AddOrderProcessed(u.ctx, u.name)
		if err != nil {
			return job.Qaitar()
		}
		u.slogger.Debug("order processing", "action", "order_completed", "order_number", order.OrderNumber)
		loggNoti.Timestamp = time.Now()
		loggNoti.Level = "INFO"
		loggNoti.Action = fmt.Sprintf("ack the order : %s", order.OrderNumber)
		loggNoti.Message = "the order is ready"
		loggNoti.Details.NewStatus = "ready"
		defer u.rabbit.PublishNotify(u.ctx, loggNoti)
		return errors.Join(err, job.Rastau())
	}
}
