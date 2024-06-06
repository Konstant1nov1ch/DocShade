package queue_service

import (
	anonymizer_service "queue-service/providers/py-anonymizer_provider"
	"queue-service/providers/rabbitmq_provider"
	"queue-service/providers/s3_provider"
)

type QueueServiceFactory interface {
	GetService() QueueService
}

type queueServiceFactory struct {
	rabbitmq   rabbitmq_provider.RabbitMQ
	s3         s3_provider.S3
	anonymizer anonymizer_service.Anonymizer
}

func NewQueueFactory(rabbitmq rabbitmq_provider.RabbitMQ, s3 s3_provider.S3, anonymizer anonymizer_service.Anonymizer) QueueServiceFactory {
	return &queueServiceFactory{
		rabbitmq:   rabbitmq,
		s3:         s3,
		anonymizer: anonymizer,
	}
}

func (c *queueServiceFactory) GetService() QueueService {
	return newQueueService(c.rabbitmq, c.s3, c.anonymizer)
}

func newQueueService(rabbitmq rabbitmq_provider.RabbitMQ, s3 s3_provider.S3, anonymizer anonymizer_service.Anonymizer) QueueService {
	return &queueService{
		rabbitmq:   rabbitmq,
		s3Service:  s3,
		anonymizer: anonymizer,
	}
}
