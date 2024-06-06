package health

// import (
// 	"bytes"
// 	"context"
// 	rest_service "document-upload-service/usecases/upload_service"
// 	mockUsecase "document-upload-service/usecases/upload_service/mocks"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	httpUtils "gitlab.com/docshade/common/http"

// 	"github.com/labstack/echo/v4"
// 	"github.com/stretchr/testify/assert"

// 	mockProviders "document-upload-service/providers/mocks"
// )

// type testingObject struct {
// 	*testing.T
// }

// func TestDo_SuccessfulHealth(t *testing.T) {
// 	// Создание мока сервиса склада
// 	testObject := &testingObject{T: t}
// 	mockRestService := mockUsecase.NewRestService(testObject)
// 	mockRestService.On("GetHealth", context.Background(), rest_service.HealthDtoIn{
// 		Message: "Artem",
// 	}).Return(
// 		rest_service.HealthDtoOut{Message: "hello Artem"}, nil)

// 	mockPr := mockProviders.NewExecutorProviders(testObject)
// 	mockFactory := mockUsecase.NewRestServiceFactory(testObject)
// 	mockPr.On("GetRestServiceFactory").Return(mockFactory)
// 	mockFactory.On("GetService").Return(mockRestService)

// 	handler := NewHealth(httpUtils.GetMethod, "/v1/health", mockPr)

// 	// Создание тестового контекста Echo
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodGet, "/v1/health", bytes.NewBufferString(`{"message": "Artem"}`))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	// Выполнение запроса
// 	err := handler.Do(c)

// 	// Проверка, что запрос завершился успешно
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, rec.Code)
// }
