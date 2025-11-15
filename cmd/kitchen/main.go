package kitchen

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	psql "pizza/internal/adapters/psql/kitchen"
	rabb "pizza/internal/adapters/rabbit/kitchen"
	"pizza/internal/config"
	"pizza/internal/services"
	"pizza/pkg"
	"strings"
	"syscall"
)

var (
	workerName = flag.String("worker-name", "", "name of worker")
	heartbeat  = flag.Uint("heartbeat-interval", 30, "Interval (seconds) between heartbeats.")
	prefetch   = flag.Uint("prefetch", 1, "RabbitMQ prefetch count, limiting how many messages the worker receives at once.")
	orderTypes = []string{}
)

func flagger() error {
	orderType := flag.String("order-types", "", "type")
	flag.Parse()
	*workerName = strings.TrimSpace(*workerName)
	if *workerName == "" {
		return fmt.Errorf("invalid name worker")
	}

	if *orderType == "" {
		*orderType = "dinein,takeout,delivery"
	}
	orderTypesMap := make(map[string]struct{})
	for v := range strings.SplitSeq(*orderType, ",") {
		if v == "dinein" || v == "takeout" || v == "delivery" {
			_, ok := orderTypesMap[v]
			if ok {
				return fmt.Errorf("duplicated order type: %s", v)
			}
			orderTypesMap[v] = struct{}{}
			orderTypes = append(orderTypes, v)
		} else {
			return fmt.Errorf("invalid order type: %s", v)
		}
	}
	return nil
}

func Main() {
	logger := pkg.CustomSlog().With("service", "kitchen-worker")
	err := flagger()
	if err != nil {
		logger.Error("invalid flag or flag not seted", "error", err)
		os.Exit(1)
	}
	fmt.Println(*workerName, orderTypes, *heartbeat, *prefetch)
	hostName, err := os.Hostname()
	if err != nil {
		logger.Error("cannot get host name")
		return
	}
	dbcfg, err := config.GetDBConfig()
	if err != nil {
		logger.Error("cannot get db config", "error", err)
		return
	}

	rabbitCfg, err := config.GetRabbitMQConfig()
	if err != nil {
		logger.Error("cannot get rabbit config", "error", err)
		return
	}

	db, err := psql.NewOrderDB(context.Background(), logger, dbcfg)
	if err != nil {
		logger.Error("cannot connect to db", "error", err)
		return
	}
	defer db.CloseDB()
	rabbit, err := rabb.NewKitchenRabbit(rabbitCfg, logger, int(*prefetch), orderTypes)
	if err != nil {
		logger.Error("cannot connect to rabbit", "error", err)
		return
	}
	defer rabbit.CloseRabbit()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	service, err := services.NewKitchenService(ctx, logger, rabbit, db, *workerName, hostName, orderTypes)
	if err != nil {
		logger.Error("cannot create service", "error", err)
		return
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		service.StartWork()
	}()

	<-quit
	err = db.UpdateToOffline(ctx, *workerName)
	if err != nil {
		logger.Error("sql", "error", err)
	}
	cancel()
	logger.Info("📦 Shutting down server...")
	logger.Info("✅ Server exited properly")
}
