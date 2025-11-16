package route

import (
	"net/http"

	"pizza/internal/adapters/server/order/route/handler"
	"pizza/internal/ports"
)

func NewRoute(use ports.OrderUseCase) http.Handler {
	mux := http.NewServeMux()
	hand := handler.NewHandler(use)
	mux.HandleFunc("POST /orders", hand.CreateOrder)
	return mux
}
