package kitchen

import (
	"context"
)

func (p *psql) UpdateStatusOrder(ctx context.Context, orderNumber, status, processedBy string) error {
	// Обновляем статус, processed_by и updated_at
	tx, err := p.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var id uint
	err = tx.QueryRow(ctx, `
        UPDATE orders
        SET status = $1,
            processed_by = $2,
            updated_at = NOW()
        WHERE number = $3
		RETURNING id;
    `, status, processedBy, orderNumber).Scan(&id)
	if err != nil {
		return err
	}
	const queryStatusLog = `
		INSERT INTO order_status_log (order_id, status, changed_by)
		VALUES ($1, $2, $3)`
	_, err = tx.Exec(ctx, queryStatusLog, id, status, processedBy)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
