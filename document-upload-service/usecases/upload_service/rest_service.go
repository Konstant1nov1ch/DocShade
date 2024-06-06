package rest_service

import (
	"context"
	"document-upload-service/providers/rabbitmq_provider"
	"document-upload-service/providers/s3_provider"
	"encoding/json"
	"errors"
)

type RestService interface {
	GetHealth(ctx context.Context, data HealthDtoIn) (HealthDtoOut, error)
	UploadDocument(ctx context.Context, sessionID, documentID, originalFileName string, fileData []byte) error
}

type restService struct {
	s3Service s3_provider.S3
	rabbitmq  rabbitmq_provider.RabbitMQ
}

// NewRestService конструктор сервиса работы с файлами
func NewRestService(s3Service s3_provider.S3, rabbitmq rabbitmq_provider.RabbitMQ) RestService {
	return &restService{
		s3Service: s3Service,
		rabbitmq:  rabbitmq,
	}
}

func (r *restService) GetHealth(ctx context.Context, data HealthDtoIn) (HealthDtoOut, error) {
	return HealthDtoOut{Message: "hello " + data.Message}, nil
}

func (r *restService) UploadDocument(ctx context.Context, sessionID, documentID, originalFileName string, fileData []byte) error {
	// Определяем путь в S3
	bucket := "preprocessing"
	objectName := documentID + ".pdf"

	// Загрузка файла в S3
	err := r.s3Service.Put(ctx, objectName, bucket, fileData, nil)
	if err != nil {
		return errors.New("failed to upload file to S3: " + err.Error())
	}

	// Создание сообщения для RabbitMQ
	messageBody := map[string]interface{}{
		"session_id":         sessionID,
		"document_id":        documentID,
		"s3_path":            bucket + "/" + objectName + ".pdf",
		"original_file_name": originalFileName,
	}
	message, err := json.Marshal(messageBody)
	if err != nil {
		return errors.New("failed to create message: " + err.Error())
	}

	// Публикация сообщения в RabbitMQ
	err = r.rabbitmq.PublishMessage(ctx, "document-exchange", "in-routing-key", message)
	if err != nil {
		return errors.New("failed to publish message to RabbitMQ: " + err.Error())
	}

	return nil
}
