package tasks

import (
	"context"
	"encoding/json"
	"log"
	"notification-service/usecases/notifi_service"
	"sync"
	"time"

	"gitlab.com/docshade/common/http"
)

type WorkerPool struct {
	notifiService notifi_service.NotifiService
	wsServer      *http.WebSocketServer
	jobs          chan notifi_service.DocumentMessage
	wg            sync.WaitGroup
	mu            sync.Mutex
	activeWorkers int
	maxWorkers    int
}

func NewWorkerPool(notifiService notifi_service.NotifiService, wsServer *http.WebSocketServer, maxWorkers int) *WorkerPool {
	pool := &WorkerPool{
		notifiService: notifiService,
		wsServer:      wsServer,
		jobs:          make(chan notifi_service.DocumentMessage),
		maxWorkers:    maxWorkers,
	}

	return pool
}

func (p *WorkerPool) worker(ctx context.Context) {
	defer func() {
		p.mu.Lock()
		p.activeWorkers--
		p.mu.Unlock()
		p.wg.Done()
	}()

	for {
		select {
		case msg := <-p.jobs:
			err := p.notifiService.ProcessDocumentMessage(ctx, msg)
			if err != nil {
				log.Printf("Failed to process message: %v", err)
				continue
			}

			// Используйте отдельный контекст для генерации временной ссылки
			genCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			downloadLink, err := p.notifiService.GeneratePresignedURL(genCtx, msg.DocumentID+".pdf", 15*time.Minute)
			if err != nil {
				log.Printf("Failed to generate presigned URL: %v", err)
				continue
			}

			// Отправка уведомления клиенту через WebSocket
			notification := map[string]interface{}{
				"session_id":        msg.SessionID,
				"status":            msg.Status,
				"download_link":     downloadLink,
				"original_filename": msg.OriginalFileName,
			}
			if msg.Status != "ok" {
				notification["status"] = msg.Status
			}
			notificationBytes, _ := json.Marshal(notification)
			log.Printf("Sending message to session %s: %s", msg.SessionID, string(notificationBytes))
			p.wsServer.SendMessageToClient(msg.SessionID, notificationBytes)
		case <-ctx.Done():
			log.Printf("Worker timed out and is shutting down")
			return
		}
	}
}

func (p *WorkerPool) AddJob(msg notifi_service.DocumentMessage) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Запуск новой горутины, если текущих горутин недостаточно
	if p.activeWorkers < p.maxWorkers {
		p.wg.Add(1)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()
		go p.worker(ctx)
		p.activeWorkers++
		log.Printf("Started new worker, active workers: %d", p.activeWorkers)
	}
	p.jobs <- msg
	log.Printf("Job added for session %s", msg.SessionID)
}

func (p *WorkerPool) Wait() {
	close(p.jobs)
	p.wg.Wait()
}
