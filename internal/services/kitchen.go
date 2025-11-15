package services

import (
	"context"
	"fmt"
	"log/slog"
	"pizza/internal/ports"
)

type kitchen struct {
	ctx     context.Context
	slogger *slog.Logger
	rabbit  ports.KitchenRabbit
	db      ports.KitchenPsql
	types   map[string]struct{}
}

func NewKitchenService(ctx context.Context, slogger *slog.Logger, rabbit ports.KitchenRabbit, db ports.KitchenPsql, name string, types []string) (ports.KitchenService, error) {
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
		ctx:     ctx,
		slogger: slogger,
		rabbit:  rabbit,
		db:      db,
		types:   typesMap,
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
			fmt.Println("keldi")
			go u.worker(job)
		}
	}
}

func (u *kitchen) worker(job ports.QatJoldama) error {
	// defer job.Qaitar()
	order, err := job.GiveBody()
	if err != nil {
		return job.Joi()
	}
	fmt.Println(order)
	return job.Rastau()
}
