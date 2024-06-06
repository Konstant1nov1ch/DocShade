package middleware

import (
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Middleware interface {
	// GetCorsConfigMiddleware корс конфигурация
	GetCorsConfigMiddleware() echo.MiddlewareFunc
	// GetGlobalMiddlewares получить набор глобальных на уровне приложения промежуточных функций
	GetGlobalMiddlewares() []echo.MiddlewareFunc
}

type BaseMW struct {
	CorsConfigMiddleware    echo.MiddlewareFunc
	JwtValidationMiddleware echo.MiddlewareFunc
}

// NewBaseMiddleware базовый набор промежуточный функций
func NewBaseMiddleware() *BaseMW {
	return &BaseMW{
		CorsConfigMiddleware: CorsConfigMiddleware(),
	}
}

func CorsConfigMiddleware() echo.MiddlewareFunc {
	defaultConfig := echomw.CORSConfig{
		Skipper:      echomw.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}

	return echomw.CORSWithConfig(defaultConfig)
}

func (mw *BaseMW) GetCorsConfigMiddleware() echo.MiddlewareFunc {
	return mw.CorsConfigMiddleware
}

func (mw *BaseMW) GetGlobalMiddlewares() []echo.MiddlewareFunc {
	glob := make([]echo.MiddlewareFunc, 0, 3)
	glob = append(glob,
		mw.GetCorsConfigMiddleware(),
	)

	return glob
}
