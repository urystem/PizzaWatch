package kitchen

import (
	"context"
	"fmt"
	"pizza/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type psql struct {
	*pgxpool.Pool
}

func NewOrderDB(ctx context.Context, cfg config.CfgDBInter) (any, error) {
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
	return &psql{db}, nil
}
