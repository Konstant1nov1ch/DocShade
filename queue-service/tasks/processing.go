package tasks

import (
	"context"
	"log"
	"queue-service/usecases/queue_service"
)

func StartQueueListener(ctx context.Context, queueService queue_service.QueueService) {
	err := queueService.ConsumeMessages(ctx, "in_queue", queueService.ProcessDocumentMessage)
	if err != nil {
		log.Fatalf("Failed to start queue listener: %v", err)
	}
}
