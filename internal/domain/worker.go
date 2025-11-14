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
