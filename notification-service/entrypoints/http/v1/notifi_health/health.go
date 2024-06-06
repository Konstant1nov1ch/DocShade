package notifi_health

import (
	"encoding/json"
	"net/http"
	notifi_service "notification-service/usecases/notifi_service"

	httpUtils "gitlab.com/docshade/common/http"

	"gitlab.com/docshade/common/core"

	"github.com/labstack/echo/v4"
)

const (
	Route  = "/v1/notifi_health"
	Method = httpUtils.GetMethod
)

type providerHealth interface {
	GetNotifiServiceFactory() notifi_service.NotifiServiceFactory
}

type health struct {
	method    httpUtils.Methods
	route     string
	providers providerHealth
}

// NewHealth get new object
func NewHealth(
	method httpUtils.Methods,
	route string,
	providers providerHealth,
) core.Handler {
	return &health{
		method:    method,
		route:     route,
		providers: providers,
	}
}

// GetMethod Get handler method
func (h *health) GetMethod() httpUtils.Methods {
	return h.method
}

// GetRoute Get handler route
func (h *health) GetRoute() string {
	return h.route
}

// Do метод, который вызывается при обращении к ручке
// @Summary     Проверить жизнеспособность сервиса
// @Produce      json
// @Param requestBody body DtoIn true "Тело запроса"
// @Success      200
// @Router       /v1/health [get]
func (h *health) Do(ctx echo.Context) error {
	body := ctx.Request().Body
	var data DtoIn
	err := json.NewDecoder(body).Decode(&data)

	service := h.providers.GetNotifiServiceFactory().GetService()

	response, err := service.GetHealth(ctx.Request().Context(), notifi_service.HealthDtoIn{Message: data.Message})

	if err != nil {
		return httpUtils.ReturnInternalError(ctx, err, "")
	}

	return ctx.JSON(http.StatusOK, prepareResponse(response))

}

func prepareResponse(data notifi_service.HealthDtoOut) DtoOut {

	return DtoOut{Message: data.Message}
}
