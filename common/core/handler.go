package core

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/docshade/common/http"
)

type Handler interface {
	// GetMethod получить название метода ручки
	GetMethod() http.Methods
	// GetRoute получить путь ручки
	GetRoute() string
	// Do метод, который вызывается при обращении к ручке
	Do(ctx echo.Context) error
}
