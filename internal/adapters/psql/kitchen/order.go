package kitchen

import (
	"context"
	"fmt"
)

func (p *psql) UpdateStatusOrder(ctx context.Context, orderNumber, status, processedBy string) error {
	// Обновляем статус, processed_by и updated_at
	cmdTag, err := p.Exec(ctx, `
        UPDATE orders
        SET status = $1,
            processed_by = $2,
            updated_at = NOW()
        WHERE number = $3
    `, status, processedBy, orderNumber)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		// Если ни одной строки не обновилось — такого заказа нет
		return fmt.Errorf("order %s not found", orderNumber)
	}
	
	return nil
}
