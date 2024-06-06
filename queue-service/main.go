package main

import (
	"context"
	"os"
	"os/signal"
	"queue-service/entrypoints/http/v1/queue_health"
	dataproviders "queue-service/providers"
	"queue-service/tasks"
	"syscall"

	"gitlab.com/docshade/common/core"
	"gitlab.com/docshade/common/http/middleware"
	logger "gitlab.com/docshade/common/log"

	"github.com/labstack/echo/v4"
)

func main() {
	config := core.NewConfig("test")
	err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Error occurred while loading config: %+v, config: %+v", err, config)
	}
	providers, err := dataproviders.NewProviders(config)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	mw := middleware.NewBaseMiddleware()
	service := echo.New()

	microservice := core.NewMicroservice(config, service, mw.GetGlobalMiddlewares())
	logger.InitLog(config.GetLogConfig())

	addRoutes(config, providers)

	// Start the microservice
	go func() {
		microservice.Run()
	}()

	// Start the queue listener task
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue_service := providers.GetQueueServiceFactory().GetService()

	go tasks.StartQueueListener(ctx, queue_service)

	// Wait for interrupt signal to gracefully shutdown the service
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()
}

func addRoutes(config core.Config, providers dataproviders.ExecutorProviders) {
	config.AddHandler(queue_health.NewHealth(queue_health.Method, queue_health.Route, providers))
}
