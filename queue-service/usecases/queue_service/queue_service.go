package queue_service

import (
	"context"
	"encoding/json"
	"errors"
	anonymizer_provider "queue-service/providers/py-anonymizer_provider"
	"queue-service/providers/rabbitmq_provider"
	"queue-service/providers/s3_provider"
)

type QueueService interface {
	GetHealth(ctx context.Context, data HealthDtoIn) (HealthDtoOut, error)
	ProcessDocumentMessage(ctx context.Context, msg rabbitmq_provider.DocumentMessage) error
	ConsumeMessages(ctx context.Context, queueName string, handler func(context.Context, rabbitmq_provider.DocumentMessage) error) error
}

type queueService struct {
	s3Service  s3_provider.S3
	rabbitmq   rabbitmq_provider.RabbitMQ
	anonymizer anonymizer_provider.Anonymizer
}

func NewQueueService(s3Service s3_provider.S3, rabbitmq rabbitmq_provider.RabbitMQ, anonymizer anonymizer_provider.Anonymizer) QueueService {
	return &queueService{
		s3Service:  s3Service,
		rabbitmq:   rabbitmq,
		anonymizer: anonymizer,
	}
}

func (r *queueService) GetHealth(ctx context.Context, data HealthDtoIn) (HealthDtoOut, error) {
	return HealthDtoOut{Message: "hello " + data.Message}, nil
}

func (r *queueService) ProcessDocumentMessage(ctx context.Context, msg rabbitmq_provider.DocumentMessage) error {
	var textMessage string
	var notificationMessage map[string]interface{}
	// Step 1: Download the file from S3
	object, err := r.s3Service.Get(ctx, msg.S3Path, msg.DocumentID+".pdf")
	if err != nil {
		return err
	}

	// Step 2: Anonymize the document
	anonymizedDocument, errByAnonim := r.anonymizer.AnonymizeDocument(object, msg.DocumentID+".pdf")
	if errByAnonim != nil {
		textMessage = "Somthing going wrong pls try another time"
	}

	// Step 3: Upload the anonymized document to S3
	destPath := "postprocessing"
	err = r.s3Service.Put(ctx, msg.DocumentID, destPath, anonymizedDocument, nil)
	if err != nil {
		return err
	}
	//object, err := r.s3Service.Get(ctx, msg.S3Path, msg.DocumentID+".pdf")
	// Step 4: Remove the original document from the preprocessing bucket
	err = r.s3Service.Remove(ctx, msg.S3Path, msg.DocumentID+".pdf")
	if err != nil {
		return err
	}

	if errByAnonim != nil {
		notificationMessage = map[string]interface{}{
			"session_id":         msg.SessionID,
			"document_id":        msg.DocumentID,
			"s3_path":            errByAnonim,
			"original_file_name": "error process file",
			"status":             textMessage,
		}
	} else {
		// Step 5: Send a notification message
		notificationMessage = map[string]interface{}{
			"session_id":         msg.SessionID,
			"document_id":        msg.DocumentID,
			"s3_path":            destPath + "/" + msg.DocumentID,
			"original_file_name": msg.OriginalFileName,
			"status":             "ok",
		}
	}
	// notificationMessage = map[string]interface{}{
	// 	"session_id":         msg.SessionID,
	// 	"document_id":        msg.DocumentID,
	// 	"s3_path":            destPath + "/" + msg.DocumentID,
	// 	"original_file_name": msg.OriginalFileName,
	// 	"status":             "ok",
	// }
	messageBytes, err := json.Marshal(notificationMessage)
	if err != nil {
		return err
	}

	err = r.rabbitmq.PublishMessage(ctx, "document-exchange", "out-routing-key", messageBytes)
	if err != nil {
		return errors.New("failed to publish message to RabbitMQ: " + err.Error())
	}

	return nil
}

func (r *queueService) ConsumeMessages(ctx context.Context, queueName string, handler func(context.Context, rabbitmq_provider.DocumentMessage) error) error {
	return r.rabbitmq.ConsumeMessages(ctx, queueName, func(msg rabbitmq_provider.DocumentMessage) error {
		return handler(ctx, msg)
	})
}
