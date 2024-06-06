package notifi_service

import (
	"notification-service/providers/rabbitmq_provider"
	"notification-service/providers/s3_provider"
)

type NotifiServiceFactory interface {
	GetService() NotifiService
}

type notifiServiceFactory struct {
	rabbitmq rabbitmq_provider.RabbitMQ
	s3       s3_provider.S3
}

func NewNotifiFactory(rabbitmq rabbitmq_provider.RabbitMQ, s3 s3_provider.S3) NotifiServiceFactory {
	return &notifiServiceFactory{
		rabbitmq: rabbitmq,
		s3:       s3,
	}
}

func (c *notifiServiceFactory) GetService() NotifiService {
	return newNotifiService(c.rabbitmq, c.s3)
}

func newNotifiService(rabbitmq rabbitmq_provider.RabbitMQ, s3 s3_provider.S3) NotifiService {
	return &notifiService{
		rabbitmq:  rabbitmq,
		s3Service: s3,
	}
}
