package tracing

import (
	"context"
	"fmt"
	"log/slog"
	"pizza/internal/config"
	"pizza/internal/ports"

	"github.com/jackc/pgx/v5/pgxpool"
)

type psql struct {
	logg *slog.Logger
	*pgxpool.Pool
}

func NewOrderDB(ctx context.Context, cfg config.CfgDBInter, logg *slog.Logger) (ports.TrackingSQL, error) {
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
		Pool: db,
		logg: logg,
	}, db.Ping(ctx)
}

func (pool *psql) CloseDB() {
	pool.Close()
	pool.logg.Info("db connection closed")
}
