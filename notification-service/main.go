package main

import (
	"context"
	"notification-service/entrypoints/http/v1/notifi_health"
	dataproviders "notification-service/providers"
	"notification-service/tasks"
	"os"
	"os/signal"
	"syscall"

	logger "gitlab.com/docshade/common/log"

	"gitlab.com/docshade/common/core"
	"gitlab.com/docshade/common/http"
	"gitlab.com/docshade/common/http/middleware"

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
	wsServer := http.NewWebSocketServer()

	microservice := core.NewMicroservice(config, service, mw.GetGlobalMiddlewares())
	logger.InitLog(config.GetLogConfig())

	addRoutes(config, providers, service, wsServer)

	// Start the microservice
	go func() {
		microservice.Run()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	notifi_service := providers.GetNotifiServiceFactory().GetService()

	go tasks.StartQueueListener(ctx, notifi_service, wsServer, 10)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()
}

func addRoutes(config core.Config, providers dataproviders.ExecutorProviders, e *echo.Echo, wsServer *http.WebSocketServer) {
	config.AddHandler(notifi_health.NewHealth(notifi_health.Method, notifi_health.Route, providers))
	http.RegisterWebSocketRoutes(e, wsServer)
}
