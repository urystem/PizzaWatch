package notification

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	rabb "pizza/internal/adapters/rabbit/notify"
	"pizza/internal/config"
	"pizza/internal/services"
	"pizza/pkg"
)

func Main() {
	logger := pkg.CustomSlog().With("service", "notification-subscriber")
	rabbitCfg, err := config.GetRabbitMQConfig()
	if err != nil {
		logger.Error("cannot get rabbit config", "error", err)
		return
	}
	myRab, err := rabb.NewNotifyRabbit(rabbitCfg, logger)
	if err != nil {
		logger.Error("cannot connect to rabbit", "error", err)
		return
	}
	defer myRab.CloseRabbit()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := services.NewNotiServive(ctx, logger, myRab)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		service.StartNotify()
	}()

	<-quit
	cancel()
	logger.Info("📦 Shutting down server...")
	logger.Info("✅ Server exited properly")
}
