package kitchen

import (
	"context"
	"fmt"
	"pizza/internal/config"
	"pizza/internal/ports"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type psql struct {
	*pgxpool.Pool
}

func NewOrderDB(ctx context.Context, cfg config.CfgDBInter) (ports.KitchenPsql, error) {
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

func (p *psql) CloseDB() {
	p.Close()
}

func (p *psql) CreateOrUpdateWorker(ctx context.Context, name string, types []string) error {
	const query = `
	INSERT INTO workers 
	(name, type) 
	VALUES ($1, $2) 
	ON CONFLICT (name) DO UPDATE 
	SET 
		type = EXCLUDED.type, 
		status = 'online', 
		last_seen = now() 
	WHERE workers.status != 'online'`

	aff, err := p.Exec(ctx, query, name, strings.Join(types, ","))
	if err != nil {
		return err
	}
	if aff.RowsAffected() == 0 {
		return fmt.Errorf("%s", "already online")
	}
	return nil
}

func (p *psql) AddOrderProcessed(ctx context.Context, name string) error {
	const query = `
		UPDATE workers
		SET orders_processed = orders_processed + 1,
		    last_seen = now()
		WHERE name = $1 AND status = 'online';`

	res, err := p.Exec(ctx, query, name)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s", "worker not found or he is offline")
	}
	return nil
}
