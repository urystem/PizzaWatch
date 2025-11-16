package domain

import "time"

type Worker struct {
	ID              uint      `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	Name            string    `json:"name"`
	Types           []string  `json:"types"`
	Status          string    `json:"status"`
	LastSeen        time.Time `json:"last_seen"`
	OrdersProcessed uint      `json:"orders_processed"`
}

type LogMessageKitchen struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Service   string    `json:"service"`
	Hostname  string    `json:"hostname"`
	Action    string    `json:"action"`
	Message   string    `json:"message"`
	Details   `json:"details"`
}

type Details struct {
	OrderNumber string `json:"order_number"`
	NewStatus   string `json:"new_status"`
}
