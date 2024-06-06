package providers

import (
	"context"
	"log"
	anonymizer_provider "queue-service/providers/py-anonymizer_provider"
	"queue-service/providers/rabbitmq_provider"
	"queue-service/providers/s3_provider"
	queue_service "queue-service/usecases/queue_service"

	"gitlab.com/docshade/common/core"
)

type ExecutorProviders interface {
	GetQueueServiceFactory() queue_service.QueueServiceFactory
}

type executorProviders struct {
	queueFactory queue_service.QueueServiceFactory
	rabbitmq     rabbitmq_provider.RabbitMQ
	s3           s3_provider.S3
	anonymizer   anonymizer_provider.Anonymizer
}

func (p *executorProviders) GetQueueServiceFactory() queue_service.QueueServiceFactory {
	return p.queueFactory
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

	err := s3.CreateBucket(context.Background(), "postprocessing")
	if err != nil {
		//return nil, err
		log.Println(err)
		err = nil
	}

	anonymizer := anonymizer_provider.NewAnonymizer(config.GetAnonymizerConfig())
	if err := anonymizer.InitAnonymizer(); err != nil {
		return nil, err
	}

	err = rabbitmq.CreateQueueAndBind(context.Background(), "out_queue", "document-exchange", "out-routing-key")
	if err != nil {
		log.Println("ошибка подключения к  CreateQueueAndBind", err)
		return nil, err
	}

	queueFactory := queue_service.NewQueueFactory(rabbitmq, s3, anonymizer)

	return &executorProviders{
		s3:           s3,
		queueFactory: queueFactory,
		rabbitmq:     rabbitmq,
		anonymizer:   anonymizer,
	}, nil
}
