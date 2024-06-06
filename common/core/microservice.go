package core

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	http "gitlab.com/docshade/common/http"
	"gitlab.com/docshade/common/log"
)

type Microservice interface {
	// Run запустить сервис
	Run()
	// GetService получить текущий эхо сервис
	GetService() *echo.Echo
	// DisableGlobalMiddleware отключить глобальные промежуточные функции
	DisableGlobalMiddleware() Microservice
	// DisableMiddleware отключить промежуточные функции локального скоупа, относящиеся к конкретным хэндлерам
	DisableMiddleware() Microservice
}
type microservice struct {
	config  Config
	service *echo.Echo

	disableGlobalMiddleware bool
	disableMiddleware       bool
	globalMiddlewares       []echo.MiddlewareFunc
}

// NewMicroservice новый объект сервиса
func NewMicroservice(
	config Config,
	service *echo.Echo,
	globalMiddlewares []echo.MiddlewareFunc,
) Microservice {

	return &microservice{
		service:           service,
		config:            config,
		globalMiddlewares: globalMiddlewares,
	}
}

// Run запустить сервис на исполнение
func (m *microservice) Run() {
	globalGroup := m.configureGlobalMiddlewares(m.service)
	m.addRoutes(globalGroup)
	m.addSwagger(m.service)

	m.service.Logger.Fatal(m.service.Start(":" + m.config.GetPort()))
}
func (m *microservice) GetService() *echo.Echo {
	return m.service
}

func (m *microservice) DisableGlobalMiddleware() Microservice {
	m.disableGlobalMiddleware = true

	return m
}

func (m *microservice) DisableMiddleware() Microservice {
	m.disableMiddleware = true

	return m
}

func (m *microservice) addRoutes(globalGroup *echo.Group) {
	for _, handler := range m.config.GetHandlerList() {
		middlewares := make([]echo.MiddlewareFunc, 0, 10)

		switch handler.GetMethod() {
		case http.GetMethod:
			globalGroup.GET(handler.GetRoute(), handler.Do, middlewares...)
		case http.PostMethod:
			globalGroup.POST(handler.GetRoute(), handler.Do, middlewares...)
		case http.PutMethod:
			globalGroup.PUT(handler.GetRoute(), handler.Do, middlewares...)
		case http.DeleteMethod:
			globalGroup.DELETE(handler.GetRoute(), handler.Do, middlewares...)
		case http.PatchMethod:
			globalGroup.PATCH(handler.GetRoute(), handler.Do, middlewares...)
		}
	}
}

func (m *microservice) configureGlobalMiddlewares(service *echo.Echo) *echo.Group {
	service.Use(
		echo.WrapMiddleware(log.RequestIDMiddleware),
		log.EchoLogger(),
	)

	if m.disableGlobalMiddleware {
		return nil
	}

	g := service.Group("")

	for _, mw := range m.globalMiddlewares {
		g.Use(mw)
	}

	return g
}

func (m *microservice) addSwagger(service *echo.Echo) {
	service.GET(http.SwaggerRoute, echoSwagger.WrapHandler)
}
