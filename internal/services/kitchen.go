package services

import (
	"context"
	"pizza/internal/ports"
)

type kitchen struct {
	ctx    context.Context
	rabbit ports.KitchenRabbit
	db     ports.KitchenPsql
}

func NewKitchenService(ctx context.Context, rabbit ports.KitchenRabbit, db ports.KitchenPsql) ports.KitchenService {
	return &kitchen{
		ctx:    ctx,
		rabbit: rabbit,
		db:     db,
	}
}

func (u *kitchen) StartWork() {
	jobs := u.rabbit.GiveChannel()
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
	defer job.Qaitar()

	return job.Rastau()
}
