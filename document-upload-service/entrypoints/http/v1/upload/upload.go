package upload

import (
	"context"
	rest_service "document-upload-service/usecases/upload_service"
	"net/http"

	"gitlab.com/docshade/common/core"
	httpUtils "gitlab.com/docshade/common/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	Route  = "/v1/upload"
	Method = httpUtils.PostMethod
)

type providerUpload interface {
	GetRestServiceFactory() rest_service.RestServiceFactory
}

type upload struct {
	method    httpUtils.Methods
	route     string
	providers providerUpload
}

// NewUpload get new object
func NewUpload(
	method httpUtils.Methods,
	route string,
	providers providerUpload,
) core.Handler {
	return &upload{
		method:    method,
		route:     route,
		providers: providers,
	}
}

// GetMethod Get handler method
func (h *upload) GetMethod() httpUtils.Methods {
	return h.method
}

// GetRoute Get handler route
func (h *upload) GetRoute() string {
	return h.route
}

// @Summary      Upload a PDF document
// @Description  Uploads a PDF document and processes it
// @Produce      json
// @Param        file formData file true "PDF file to upload"
// @Success      200 {object} DtoOut
// @Router       /v1/upload [post]
func (h *upload) Do(ctx echo.Context) error {
	// Получение файла из запроса
	file, err := ctx.FormFile("file")
	if err != nil {
		return httpUtils.ReturnBadRequestError(ctx, err, "Invalid file")
	}

	if file.Header.Get("Content-Type") != "application/pdf" {
		return httpUtils.ReturnBadRequestError(ctx, err, "Invalid file format. Only PDF is allowed.")
	}

	// Открытие файла
	src, err := file.Open()
	if err != nil {
		return httpUtils.ReturnInternalError(ctx, err, "Failed to open file")
	}
	defer src.Close()

	// Прочтение файла в память
	fileData := make([]byte, file.Size)
	_, err = src.Read(fileData)
	if err != nil {
		return httpUtils.ReturnInternalError(ctx, err, "Failed to read file")
	}

	// Генерация идентификаторов сессии и документа
	sessionID := uuid.New().String()
	documentID := uuid.New().String()

	// Получение сервиса
	service := h.providers.GetRestServiceFactory().GetService()

	// Загрузка файла в S3 и публикация сообщения в RabbitMQ через сервис
	err = service.UploadDocument(context.Background(), sessionID, documentID, file.Filename, fileData)
	if err != nil {
		return httpUtils.ReturnInternalError(ctx, err, "Failed to process file")
	}

	response := DtoOut{
		SessionID:  sessionID,
		DocumentID: documentID,
		Message:    "File uploaded successfully",
	}
	return ctx.JSON(http.StatusOK, response)
}
