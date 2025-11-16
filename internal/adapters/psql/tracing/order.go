package tracing

import (
	"context"
	"errors"

	"pizza/internal/domain"

	"github.com/jackc/pgx/v5"
)

func (p *psql) OrderStatusUpdate(ctx context.Context, number string) (*domain.OrderStatusUpdate, error) {
	var o domain.OrderStatusUpdate
	const query = `
    SELECT 
		number, 
		status, 
		updated_at,
        COALESCE(completed_at, updated_at + interval '10 minutes') AS estimated_completion,
        processed_by
	FROM orders
    WHERE number = $1`

	err := p.QueryRow(ctx, query, number).Scan(
		&o.OrderNumber,
		&o.CurrentStatus,
		&o.UpdatedAt,
		&o.EstimatedCompletion,
		&o.ProcessedBy,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound // заказа нет
		}
		return nil, err
	}

	return &o, nil
}

func (p *psql) GetOrderHistory(ctx context.Context, number string) ([]domain.OrderStatusEvent, error) {
	var events []domain.OrderStatusEvent

	rows, err := p.Query(ctx, `
        SELECT osl.status, osl.changed_at, osl.changed_by
        FROM order_status_log osl
        JOIN orders o ON o.id = osl.order_id
        WHERE o.number = $1
        ORDER BY osl.changed_at ASC
    `, number)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e domain.OrderStatusEvent
		if err := rows.Scan(&e.Status, &e.Timestamp, &e.ChangedBy); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return events, nil
}
