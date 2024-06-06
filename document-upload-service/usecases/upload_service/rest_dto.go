package rest_service

import "context"

type RestServiceRabbitMQ interface {
	// PublishMessage публикация сообщения в RabbitMQ
	PublishMessage(ctx context.Context, exchange, routingKey string, message []byte) error
}

type RestServices3Service interface {
	// Put загружает файл в S3
	Put(ctx context.Context, objectName, path string, objectBody []byte, metaData map[string]string) error
	// IsObjectExist проверяет, существует ли объект в S3
	IsObjectExist(ctx context.Context, path, objectName string) (bool, error)
}

// HealthDtoOut Output DTO for Health Method
type HealthDtoOut struct {
	Message string
}

// HealthDtoIn Input DTO for Health Method
type HealthDtoIn struct {
	Message string
}
