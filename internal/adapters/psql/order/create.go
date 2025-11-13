package order

import (
	"context"
	"fmt"
	"pizza/internal/domain"
	"time"
)

func (p *psql) CreateOrder(ctx context.Context, ord *domain.OrderPublish) error {
	// Начинаем транзакцию
	tx, err := p.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var san uint
	var myTime time.Time

	const numQuery = `
    SELECT COUNT(*), CURRENT_DATE
    FROM orders
    WHERE created_at::date = CURRENT_DATE;`

	err = tx.QueryRow(ctx, numQuery).Scan(&san, &myTime)
	if err != nil {
		return err
	}
	ord.OrderNumber = fmt.Sprintf("ORD_%s_%03d", myTime.Format("20060102"), san)
	// Вставка заказа
	const queryOrder = `
		INSERT INTO orders (number, customer_name, order_type, priority, total_amount, table_number, delivery_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err = tx.QueryRow(ctx, queryOrder,
		ord.OrderNumber,
		ord.CustomerName,
		ord.OrderType,
		ord.Priority,
		ord.TotalAmount,
		ord.TableNumber,  // может быть NULL
		ord.DeliveryAddr, // может быть NULL
	).Scan(&san)
	if err != nil {
		return err
	}

	const queryItem = `
		INSERT INTO order_items (order_id, name, quantity, price)
		VALUES ($1, $2, $3, $4)`

	for _, item := range ord.Items {
		_, err = tx.Exec(ctx, queryItem,
			san,
			item.Name,
			item.Quantity,
			item.Price)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
