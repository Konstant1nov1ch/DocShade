package rabbitmq_provider

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.com/docshade/common/core"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	maxRetries = 10
	retryDelay = 5 * time.Second
)

type RabbitMQ interface {
	// InitRabbitMQ инициализация соединения с RabbitMQ
	InitRabbitMQ() error
	// PublishMessage публикация сообщения в RabbitMQ
	PublishMessage(ctx context.Context, exchange, routingKey string, message []byte) error
	// CreateQueueAndBind создание очереди и привязка её к обменнику
	CreateQueueAndBind(ctx context.Context, queueName, exchange, routingKey string) error
	// CreateExchange создание обменника
	CreateExchange(ctx context.Context, exchange string) error
}

type rabbitmq struct {
	cfg core.RabbitMQConfig
	mq  *amqp.Connection
}

// NewRabbitMQ создает новый экземпляр RabbitMQ
func NewRabbitMQ(cfg core.RabbitMQConfig) RabbitMQ {
	return &rabbitmq{cfg: cfg}
}

// InitRabbitMQ инициализирует соединение с RabbitMQ
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

// PublishMessage публикует сообщение в RabbitMQ
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
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	if err != nil {
		return err
	}

	return nil
}

// CreateExchange создает обменник
func (r *rabbitmq) CreateExchange(ctx context.Context, exchange string) error {
	ch, err := r.mq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	return err
}

// CreateQueueAndBind создает очередь и привязывает её к обменнику
func (r *rabbitmq) CreateQueueAndBind(ctx context.Context, queueName, exchange, routingKey string) error {
	ch, err := r.mq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		queueName,  // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

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
