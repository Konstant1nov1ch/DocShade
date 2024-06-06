package rest_service

import (
	"document-upload-service/providers/rabbitmq_provider"
	"document-upload-service/providers/s3_provider"
)

type RestServiceFactory interface {
	GetService() RestService
}

type restServiceFactory struct {
	rabbitmq rabbitmq_provider.RabbitMQ
	s3       s3_provider.S3
}

// NewRestFactory получить новый экземпляр фабрики сервисов
func NewRestFactory(rabbitmq rabbitmq_provider.RabbitMQ, s3 s3_provider.S3) RestServiceFactory {
	return &restServiceFactory{
		rabbitmq: rabbitmq,
		s3:       s3,
	}
}

// GetService получить новых экземпляр сервиса
func (c *restServiceFactory) GetService() RestService {
	return newRestService(c.rabbitmq, c.s3)
}

func newRestService(rabbitmq rabbitmq_provider.RabbitMQ, s3 s3_provider.S3) RestService {
	return &restService{
		rabbitmq:  rabbitmq,
		s3Service: s3,
	}
}
