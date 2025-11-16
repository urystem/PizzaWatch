package tracking

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	psql "pizza/internal/adapters/psql/tracing"
	srv "pizza/internal/adapters/server/tracing"
	"pizza/internal/config"
	"pizza/internal/services"
	"pizza/pkg"
)

func Main() {
	port := flag.Uint("port", 3002, "	The HTTP port for the API.")
	flag.Parse()
	logger := pkg.CustomSlog().With("service", "tracking-service")

	dbcfg, err := config.GetDBConfig()
	if err != nil {
		logger.Error("cannot get db config", "error", err)
		return
	}
	db, err := psql.NewOrderDB(context.Background(), dbcfg, logger)
	if err != nil {
		logger.Error("cannot connect to db", "error", err)
		return
	}
	defer db.CloseDB()
	service := services.NewTrackingService(logger, db)
	server := srv.NewServer(*port, service)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.StartServer(); err != nil && err != http.ErrServerClosed {
			logger.Error("❌", " Server error:", err)
			quit <- syscall.SIGTERM
		}
	}()
	<-quit
	logger.Info("📦 Shutting down server...")
	if err := server.ShutDownServer(context.Background()); err != nil {
		logger.Error("❌", " Server forced to shutdown: %v", err)
	}
	logger.Info("✅ Server exited properly")
}
