package order_service

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	psql "pizza/internal/adapters/psql/order"
	rabbit "pizza/internal/adapters/rabbit/order"
	server "pizza/internal/adapters/server/order"
	"pizza/internal/config"
	"pizza/internal/services"
	"pizza/pkg"
	"syscall"
	"time"
)

func Main() {
	logger := pkg.CustomSlog().With("service", "order-service")
	port := flag.Uint("port", 3000, "The HTTP port for the API.")
	maxCon := flag.Uint("max-concurrent", 50, "Maximum number of concurrent orders to process.")
	flag.Parse()
	logger.Info(fmt.Sprintf("service starting on port %d", *port), "action", "start the service")

	dbCfg, err := config.GetDBConfig()
	if err != nil {
		logger.Error("cannot get db config", "error", err)
		return
	}

	rabbitCfg, err := config.GetRabbitMQConfig()
	if err != nil {
		logger.Error("cannot get rabbit config", "error", err)
		return
	}

	start := time.Now()
	db, err := psql.NewOrderDB(context.Background(), dbCfg)
	if err != nil {
		logger.Error("cannot connect to db", "error", err)
		return
	}
	defer db.CloseDB()
	logger.Info("Connected to PostgreSQL database", "action", "db_connected", "duration_ms", time.Since(start).Milliseconds())

	start = time.Now()
	rab, err := rabbit.NewOrderRabbit(rabbitCfg, logger)
	if err != nil {
		logger.Error("cannot connect to rabbitMQ", "error", err)
		return
	}
	defer rab.CloseRabbit()
	logger.Info("Connected to RabbitMQ exchange 'orders_topic'", "action", "rabbitmq_connected", "duration_ms", time.Since(start).Milliseconds())

	service := services.NewOrderService(rab, db, *maxCon)
	serv := server.NewServer(*port, service)
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := serv.StartServer(); err != nil && err != http.ErrServerClosed {
			logger.Error("❌", " Server error:", err)
			quit <- syscall.SIGTERM
		}
	}()

	<-quit

	serv.ShutDownServer(context.Background())
	logger.Info("📦 Shutting down server...")
	if err := serv.ShutDownServer(context.Background()); err != nil {
		logger.Error("❌", " Server forced to shutdown: %v", err)
	}
	logger.Info("✅ Server exited properly")
}
