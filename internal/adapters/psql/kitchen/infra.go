package kitchen

import (
	"context"
	"fmt"
	"log/slog"
	"pizza/internal/config"
	"pizza/internal/ports"

	"github.com/jackc/pgx/v5/pgxpool"
)

type psql struct {
	slogg *slog.Logger
	*pgxpool.Pool
}

func NewOrderDB(ctx context.Context, logg *slog.Logger, cfg config.CfgDBInter) (ports.KitchenPsql, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		cfg.GetUser(),
		cfg.GetPassword(),
		cfg.GetHostName(),
		cfg.GetDBPort(),
		cfg.GetDBName(),
	)
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &psql{
		slogg: logg,
		Pool:  db,
	}, nil
}

func (p *psql) CloseDB() {
	p.Close()
	p.slogg.Info("db closed")
}
