package kitchen

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (p *psql) CreateOrUpdateWorker(ctx context.Context, name string, types []string) ([]string, error) {
	tx, err := p.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var status string
	err = tx.QueryRow(ctx,
		`SELECT status FROM workers WHERE name = $1`,
		name,
	).Scan(&status)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	} else if status == "online" {
		return nil, fmt.Errorf("already online")
	}

	const queryInsert = `
	INSERT INTO workers 
	(name, type) 
	VALUES ($1, $2);`

	if status == "" {
		_, err = tx.Exec(ctx, queryInsert, name, strings.Join(types, ","))
		if err != nil {
			return nil, err
		}
		p.slogg.Info("registring the new worker", "action", "worker_registered", "name", name)
		return types, tx.Commit(ctx)
	}

	const queryUpdate = `
	UPDATE workers
	SET status = 'online',
    last_seen = now()
	WHERE name = $1 AND status = 'offline'
	RETURNING type;`

	var returnedType string
	err = p.QueryRow(ctx, queryUpdate, name).Scan(&returnedType)
	if err != nil {
		return nil, err
	}
	p.slogg.Info("the worker is already registrated", "name", name)
	return strings.Split(returnedType, ","), tx.Commit(ctx)
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

func (p *psql) UpdateToOffline(ctx context.Context, name string) error {
	// Обновляем статус и last_seen
	cmdTag, err := p.Exec(ctx, `
        UPDATE workers
        SET status = 'offline',
            last_seen = NOW()
        WHERE name = $1
    `, name)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		// Если ни одной строки не обновилось, значит такого worker нет
		return fmt.Errorf("worker %s not found", name)
	}
	return nil
}
