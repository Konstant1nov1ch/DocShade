package upload

// import (
// 	"bytes"
// 	httpUtils "common/http"
// 	rest_service "document-upload-service/usecases/upload_service"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/labstack/echo/v4"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// )

// type MockProvider struct {
// 	mock.Mock
// }

// func (m *MockProvider) GetRestServiceFactory() rest_service.RestServiceFactory {
// 	args := m.Called()
// 	return args.Get(0).(rest_service.RestServiceFactory)
// }

// func TestUploadHandler(t *testing.T) {
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/v1/upload", bytes.NewReader([]byte{}))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	mockService := new(rest_service.MockUploadService)
// 	mockProvider := new(MockProvider)
// 	mockFactory := new(rest_service.MockRestServiceFactory)

// 	mockProvider.On("GetRestServiceFactory").Return(mockFactory)
// 	mockFactory.On("GetService").Return(mockService)
// 	mockService.On("UploadDocument", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

// 	h := NewUpload(httpUtils.PostMethod, Route, mockProvider)

// 	if assert.NoError(t, h.Do(c)) {
// 		require.Equal(t, http.StatusOK, rec.Code)
// 	}
// }
