package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"pizza/internal/ports"
)

type notify struct {
	ctx  context.Context
	logg *slog.Logger
	rabb ports.NotifyRabbit
}

func NewNotiServive(ctx context.Context, logg *slog.Logger, rabb ports.NotifyRabbit) ports.NotifyService {
	return &notify{
		ctx:  ctx,
		logg: logg,
		rabb: rabb,
	}
}

func (n *notify) StartNotify() {
	jobs := n.rabb.GiveChannel()
	for {
		select {
		case <-n.ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			n.prettyPrint(job)
		}
	}
}

func (n *notify) prettyPrint(data []byte) {
	var buf bytes.Buffer
	// Попытка красиво отформатировать как JSON
	if err := json.Indent(&buf, data, "", "  "); err == nil {
		fmt.Println(buf.String())
		return
	}
}
