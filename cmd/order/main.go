package order_service

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	psql "pizza/internal/adapters/psql/order"
	rabbit "pizza/internal/adapters/rabbit/order"
	server "pizza/internal/adapters/server/order"
	"pizza/internal/config"
	"pizza/internal/services"
	"syscall"
)

func Main() {
	port := flag.Uint("port", 3000, "The HTTP port for the API.")
	maxCon := flag.Uint("max-concurrent", 50, "Maximum number of concurrent orders to process.")
	// flag.Usage=
	flag.Parse()
	dbCfg, err := config.GetDBConfig()
	if err != nil {
		slog.Error("")
		return
	}
	rabbitCfg, err := config.GetRabbitMQConfig()
	if err != nil {
		slog.Error("")
		return
	}
	db, err := psql.NewOrderDB(context.Background(), dbCfg)
	if err != nil {
		slog.Error("")
		return
	}

	rab, err := rabbit.NewOrderRabbit(rabbitCfg, *maxCon)
	if err != nil {
		slog.Error("")
		return
	}
	service := services.NewOrderService(rab, db)
	serv := server.NewServer(*port, service)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := serv.StartServer(); err != nil && err != http.ErrServerClosed {
			slog.Error("❌", " Server error:", err)
			quit <- syscall.SIGTERM
		}
	}()
	<-quit
	slog.Info("📦 Shutting down server...")
	if err := serv.ShutDownServer(context.Background()); err != nil {
		slog.Error("❌", " Server forced to shutdown: %v", err)
	}
	slog.Info("✅ Server exited properly")
}
