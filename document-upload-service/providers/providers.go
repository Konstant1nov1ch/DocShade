package providers

import (
	"context"
	"document-upload-service/providers/rabbitmq_provider"
	"document-upload-service/providers/s3_provider"
	rest_service "document-upload-service/usecases/upload_service"
	"log"

	"gitlab.com/docshade/common/core"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ExecutorProviders
type ExecutorProviders interface {
	// GetRestServiceFactory  получить фабрику для работы с логикой пользователя
	GetRestServiceFactory() rest_service.RestServiceFactory
}

type executorProviders struct {
	restFactory rest_service.RestServiceFactory
	rabbitmq    rabbitmq_provider.RabbitMQ
	s3          s3_provider.S3
}

func (p *executorProviders) GetRestServiceFactory() rest_service.RestServiceFactory {
	return p.restFactory
}

// NewProviders инициализация провайдеров
func NewProviders(config core.Config) (ExecutorProviders, error) {
	s3 := s3_provider.NewS3(config.GetS3Config())
	if err := s3.InitS3(); err != nil {
		log.Println("ошибка подключения к s3 ", err)
		return nil, err
	}

	rabbitmq := rabbitmq_provider.NewRabbitMQ(config.GetRabbitMQConfig())
	if err := rabbitmq.InitRabbitMQ(); err != nil {
		log.Println("ошибка подключения к  rabbitmq", err)
		return nil, err
	}
	err := s3.CreateBucket(context.Background(), "preprocessing")
	if err != nil {
		//return nil, err
		log.Println(err)
		err = nil
	}

	// Создание обменника
	err = rabbitmq.CreateExchange(context.Background(), "document-exchange")
	if err != nil {
		log.Println("ошибка подключения к  CreateExchange", err)
		return nil, err
	}

	// Создание очереди и привязка её к обменнику
	err = rabbitmq.CreateQueueAndBind(context.Background(), "in_queue", "document-exchange", "in-routing-key")
	if err != nil {
		log.Println("ошибка подключения к  CreateQueueAndBind", err)
		return nil, err
	}

	restFactory := rest_service.NewRestFactory(rabbitmq, s3)

	return &executorProviders{
		s3:          s3,
		restFactory: restFactory,
		rabbitmq:    rabbitmq,
	}, nil
}
