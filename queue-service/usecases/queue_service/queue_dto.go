package queue_service

import "context"

type QueueServiceRabbitMQ interface {
	PublishMessage(ctx context.Context, exchange, routingKey string, message []byte) error
	CreateQueueAndBind(ctx context.Context, queueName, exchange, routingKey string) error
	CreateExchange(ctx context.Context, exchange string) error
	ConsumeMessages(ctx context.Context, queueName string, handler func(DocumentMessage) error) error
}

type QueueServices3Service interface {
	Put(ctx context.Context, objectName, path string, objectBody []byte, metaData map[string]string) error
	IsObjectExist(ctx context.Context, path, objectName string) (bool, error)
	Remove(ctx context.Context, objectName, path string) error
	Move(ctx context.Context, objectName, srcPath, destPath, newDirName string) (string, error)
}

type QueueAnonymizerService interface {
	AnonymizeDocument(document []byte) ([]byte, error)
}

// HealthDtoOut Output DTO for Health Method
type HealthDtoOut struct {
	Message string
}

type DocumentMessage struct {
	DocumentID       string `json:"document_id"`
	OriginalFileName string `json:"original_file_name"`
	S3Path           string `json:"s3_path"`
	SessionID        string `json:"session_id"`
}

// HealthDtoIn Input DTO for Health Method
type HealthDtoIn struct {
	Message string
}
