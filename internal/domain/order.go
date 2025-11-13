package domain

import "time"

// POST
type Order struct {
	CustomerName string  `json:"customer_name"`
	OrderType    string  `json:"order_type"`
	DeliveryAddr *string `json:"delivery_address,omitempty"`
	TableNumber  *uint   `json:"table_number,omitempty"`
	Items        []Item  `json:"items"`
}

type Item struct {
	Name     string  `json:"name"`
	Quantity uint    `json:"quantity"`
	Price    float64 `json:"price"`
}

// Response
type OrderStatus struct {
	OrderNumber string  `json:"order_number"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`
}
type OrderPublish struct {
	OrderNumber string `json:"order_number"`
	Order
	Priority    uint8   `json:"priority"`
	TotalAmount float64 `json:"total_amount"`
}

type OrderNotification struct {
	OrderNumber string    `json:"order_number"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	Timestamp   time.Time `json:"timestamp"`
	Priority    uint8     `json:"priority"`
}
