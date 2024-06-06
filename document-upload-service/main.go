package main

import (
	"document-upload-service/entrypoints/http/v1/health"
	"document-upload-service/entrypoints/http/v1/upload"
	dataproviders "document-upload-service/providers"

	"gitlab.com/docshade/common/core"
	"gitlab.com/docshade/common/http/middleware"
	logger "gitlab.com/docshade/common/log"

	_ "document-upload-service/docs"

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

	microservice.Run()
}

func addRoutes(
	config core.Config,
	providers dataproviders.ExecutorProviders) {

	config.AddHandler(health.NewHealth(health.Method, health.Route, providers)).
		AddHandler(upload.NewUpload(upload.Method, upload.Route, providers))

}
