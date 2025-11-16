package tracing

import (
	"context"

	"pizza/internal/domain"
)

func (p *psql) GetWorkers(ctx context.Context, heartbeatInterval uint) ([]domain.WorkerStatus, error) {
	// Порог для offline в секундах (например, 2 * heartbeatInterval)
	// const heartbeatInterval = 10 // сек
	offlineThreshold := 2 * heartbeatInterval

	rows, err := p.Query(ctx, `
        SELECT name, 
               CASE 
                 WHEN EXTRACT(EPOCH FROM (NOW() - last_seen)) > $1 THEN 'offline'
                 ELSE status
               END AS status,
               orders_processed,
               last_seen
        FROM workers;`, offlineThreshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workers []domain.WorkerStatus

	for rows.Next() {
		var w domain.WorkerStatus
		if err := rows.Scan(&w.WorkerName, &w.Status, &w.OrdersProcessed, &w.LastSeen); err != nil {
			return nil, err
		}
		workers = append(workers, w)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return workers, nil
}
