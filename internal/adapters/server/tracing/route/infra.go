package route

import (
	"net/http"

	"pizza/internal/adapters/server/tracing/route/handler"
	"pizza/internal/ports"
)

func NewRoute(use ports.TrackingUse) http.Handler {
	mux := http.NewServeMux()
	hand := handler.NewHandler(use)
	mux.HandleFunc("GET /orders/{order_number}/status", hand.GetOrderStatus)
	mux.HandleFunc("GET /orders/{order_number}/history", hand.GetOrderHistory)
	mux.HandleFunc("GET /workers/status", hand.GetWorkers)
	return mux
}
