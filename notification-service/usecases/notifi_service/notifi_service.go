package notifi_service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notification-service/providers/rabbitmq_provider"
	"notification-service/providers/s3_provider"
	"time"
)

type NotifiService interface {
	GetHealth(ctx context.Context, data HealthDtoIn) (HealthDtoOut, error)
	ProcessDocumentMessage(ctx context.Context, msg DocumentMessage) error
	ConsumeMessages(ctx context.Context, queueName string, handler func(context.Context, DocumentMessage) error) error
	GeneratePresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
	GetFileData(ctx context.Context, documentID string) ([]byte, error)
}

type notifiService struct {
	s3Service s3_provider.S3
	rabbitmq  rabbitmq_provider.RabbitMQ
}

func NewNotifiService(s3Service s3_provider.S3, rabbitmq rabbitmq_provider.RabbitMQ) NotifiService {
	return &notifiService{
		s3Service: s3Service,
		rabbitmq:  rabbitmq,
	}
}

func (r *notifiService) GetFileData(ctx context.Context, documentID string) ([]byte, error) {
	return r.s3Service.Get(ctx, documentID)
}

func (r *notifiService) GeneratePresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	return r.s3Service.GeneratePresignedURL(ctx, objectName, expiry)
}

func (r *notifiService) GetHealth(ctx context.Context, data HealthDtoIn) (HealthDtoOut, error) {
	return HealthDtoOut{Message: "hello " + data.Message}, nil
}

func (r *notifiService) ProcessDocumentMessage(ctx context.Context, msg DocumentMessage) error {
	var notificationMessage map[string]interface{}

	if msg.Status == "ok" {
		// Генерация ссылки для скачивания файла
		downloadLink := fmt.Sprintf("http://localhost:9000/%s", msg.S3Path)
		notificationMessage = map[string]interface{}{
			"session_id":    msg.SessionID,
			"status":        "ok",
			"download_link": downloadLink,
		}
	} else {
		// Сообщение об ошибке
		notificationMessage = map[string]interface{}{
			"session_id": msg.SessionID,
			"status":     "error",
		}
	}

	messageBytes, err := json.Marshal(notificationMessage)
	if err != nil {
		return err
	}

	log.Println(messageBytes, notificationMessage)

	return nil
}

func (r *notifiService) ConsumeMessages(ctx context.Context, queueName string, handler func(context.Context, DocumentMessage) error) error {
	return r.rabbitmq.ConsumeMessages(ctx, queueName, func(msg rabbitmq_provider.DocumentMessage) error {
		documentMsg := DocumentMessage{
			DocumentID:       msg.DocumentID,
			OriginalFileName: msg.OriginalFileName,
			S3Path:           msg.S3Path,
			SessionID:        msg.SessionID,
			Status:           msg.Status,
		}
		return handler(ctx, documentMsg)
	})
}
