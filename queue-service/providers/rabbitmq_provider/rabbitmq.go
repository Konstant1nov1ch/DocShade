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
	maxRetries = 10
	retryDelay = 5 * time.Second
)

type RabbitMQ interface {
	InitRabbitMQ() error
	PublishMessage(ctx context.Context, exchange, routingKey string, message []byte) error
	CreateQueueAndBind(ctx context.Context, queueName, exchange, routingKey string) error
	CreateExchange(ctx context.Context, exchange string) error
	ConsumeMessages(ctx context.Context, queueName string, handler func(DocumentMessage) error) error
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

func (r *rabbitmq) PublishMessage(ctx context.Context, exchange, routingKey string, message []byte) error {
	ch, err := r.mq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = r.CreateExchange(ctx, exchange)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	if err != nil {
		return err
	}

	return nil
}

func (r *rabbitmq) CreateExchange(ctx context.Context, exchange string) error {
	ch, err := r.mq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (r *rabbitmq) CreateQueueAndBind(ctx context.Context, queueName, exchange, routingKey string) error {
	var err error
	var ch *amqp.Channel
	for i := 0; i < maxRetries; i++ {
		ch, err = r.mq.Channel()
		if err == nil {
			defer ch.Close()
			_, err = ch.QueueDeclare(
				queueName,
				true,  // durable
				false, // autoDelete
				false, // exclusive
				false, // noWait
				nil,   // args
			)
			if err != nil {
				ch.Close()
				continue
			}

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

// Helper function to create the message
func CreateMessage(sessionID, documentID, s3Path, originalFileName string, metadata map[string]string) ([]byte, error) {
	message := map[string]interface{}{
		"session_id":         sessionID,
		"document_id":        documentID,
		"s3_path":            s3Path,
		"original_file_name": originalFileName,
		"metadata":           metadata,
	}

	return json.Marshal(message)
}
