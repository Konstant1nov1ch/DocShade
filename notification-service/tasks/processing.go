package tasks

import (
	"context"
	"log"
	"notification-service/usecases/notifi_service"

	"gitlab.com/docshade/common/http"
)

func StartQueueListener(ctx context.Context, notifiService notifi_service.NotifiService, wsServer *http.WebSocketServer, maxWorkers int) {
	pool := NewWorkerPool(notifiService, wsServer, maxWorkers)

	err := notifiService.ConsumeMessages(ctx, "out_queue", func(ctx context.Context, msg notifi_service.DocumentMessage) error {
		log.Printf("Message received for session %s", msg.SessionID)
		pool.AddJob(msg)
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to start queue listener: %v", err)
	}

	<-ctx.Done()
	pool.Wait()
}
