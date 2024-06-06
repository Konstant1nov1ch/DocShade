package providers

import (
	"context"
	"log"
	"notification-service/providers/rabbitmq_provider"
	"notification-service/providers/s3_provider"
	notifi_service "notification-service/usecases/notifi_service"

	"gitlab.com/docshade/common/core"
)

type ExecutorProviders interface {
	GetNotifiServiceFactory() notifi_service.NotifiServiceFactory
}

type executorProviders struct {
	notifiFactory notifi_service.NotifiServiceFactory
	rabbitmq      rabbitmq_provider.RabbitMQ
	s3            s3_provider.S3
}

func (p *executorProviders) GetNotifiServiceFactory() notifi_service.NotifiServiceFactory {
	return p.notifiFactory
}

// NewProviders инициализация провайдеров
func NewProviders(config core.Config) (ExecutorProviders, error) {
	s3 := s3_provider.NewS3(config.GetS3Config())
	if err := s3.InitS3(); err != nil {
		log.Println("ошибка подключения к s3 ", err)
		return nil, err
	}

	rabbitmq := rabbitmq_provider.InitRabbitMQ(config.GetRabbitMQConfig())
	if err := rabbitmq.InitRabbitMQ(); err != nil {
		log.Println("ошибка подключения к  rabbitmq", err)
		return nil, err
	}
	//and here
	err := rabbitmq.BindQueue(context.Background(), "out_queue", "document-exchange", "out-routing-key")
	if err != nil {
		log.Println("ошибка подключения к  CreateQueueAndBind", err)
		return nil, err
	}

	notifiFactory := notifi_service.NewNotifiFactory(rabbitmq, s3)

	return &executorProviders{
		s3:            s3,
		notifiFactory: notifiFactory,
		rabbitmq:      rabbitmq,
	}, nil
}
