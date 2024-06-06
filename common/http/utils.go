package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ReturnInternalError вернуть детализированную ошибку 500
func ReturnInternalError(ctx echo.Context, err error, detail string) error {
	return ctx.JSON(http.StatusInternalServerError, ErrorHttp{
		ErrorText: fmt.Sprintf("%s", err),
		Details:   []string{detail},
	})

}

// ReturnBadRequestError Вернуть ошибку плохого запроса 400
func ReturnBadRequestError(ctx echo.Context, err error, detail string) error {
	return ctx.JSON(http.StatusBadRequest, ErrorHttp{
		ErrorText: fmt.Sprintf("%s", err),
		Details:   []string{detail},
	})
}
