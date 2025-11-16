package order

import (
	"context"
	"fmt"
	"net/http"

	"pizza/internal/adapters/server/order/route"
	"pizza/internal/ports"
)

type server struct {
	http.Server
}

func NewServer(port uint, use ports.OrderUseCase) ports.ServerInter {
	return &server{
		http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: route.NewRoute(use),
		},
	}
}

func (s *server) StartServer() error {
	return s.ListenAndServe()
}

func (s *server) ShutDownServer(ctx context.Context) error {
	return s.Shutdown(ctx)
}
