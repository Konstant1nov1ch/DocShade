package rabbitmq_provider

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"gitlab.com/docshade/common/core"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	maxRetries = 11
	retryDelay = 5 * time.Second
)

type RabbitMQ interface {
	InitRabbitMQ() error
	ConsumeMessages(ctx context.Context, queueName string, handler func(DocumentMessage) error) error
	BindQueue(ctx context.Context, queueName, exchange, routingKey string) error
}

type rabbitmq struct {
	cfg core.RabbitMQConfig
	mq  *amqp.Connection
}

func InitRabbitMQ(cfg core.RabbitMQConfig) RabbitMQ {
	return &rabbitmq{cfg: cfg}
}

func (r *rabbitmq) InitRabbitMQ() error {
	var err error
	for i := 0; i < maxRetries; i++ {
		r.mq, err = amqp.Dial(r.cfg.URI)
		if err == nil {
			break
		}
		time.Sleep(retryDelay)
	}
	return nil
}

func (r *rabbitmq) BindQueue(ctx context.Context, queueName, exchange, routingKey string) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		ch, err := r.mq.Channel()
		if err == nil {
			defer ch.Close()
			err = ch.QueueBind(
				queueName,
				routingKey,
				exchange,
				false, // noWait
				nil,   // args
			)
			if err != nil {
				ch.Close()
				continue
			}
			return nil
		}
		time.Sleep(retryDelay)
	}
	return err
}

func (r *rabbitmq) ConsumeMessages(ctx context.Context, queueName string, handler func(DocumentMessage) error) error {

	//Todo error
	ch, err := r.mq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var msg DocumentMessage
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			err = handler(msg)
			if err != nil {
				log.Printf("Failed to process message: %v", err)
			}
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-ctx.Done()

	return nil
}
