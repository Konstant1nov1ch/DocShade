package http

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketServer struct {
	clients       map[string]*websocket.Conn
	messageQueues map[string][][]byte
	mu            sync.Mutex
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients:       make(map[string]*websocket.Conn),
		messageQueues: make(map[string][][]byte),
	}
}

func (server *WebSocketServer) handleWebSocket(c echo.Context) error {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	clientID := c.Param("id")
	log.Printf("Client %s connected", clientID)

	server.mu.Lock()
	server.clients[clientID] = conn
	messageQueue := server.messageQueues[clientID]
	delete(server.messageQueues, clientID)
	server.mu.Unlock()

	// Send queued messages
	for _, msg := range messageQueue {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("Failed to send queued message to client %s: %v", clientID, err)
			break
		}
		log.Printf("Queued message sent to client %s", clientID)
	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			server.mu.Lock()
			delete(server.clients, clientID)
			server.mu.Unlock()
			log.Printf("Client %s disconnected", clientID)
			break
		}
	}
	return nil
}

func (server *WebSocketServer) SendMessageToClient(clientID string, message []byte) error {
	server.mu.Lock()
	defer server.mu.Unlock()

	if conn, ok := server.clients[clientID]; ok {
		log.Printf("Sending message to connected client %s", clientID)
		return conn.WriteMessage(websocket.TextMessage, message)
	}

	log.Printf("Queueing message for client %s", clientID)
	server.messageQueues[clientID] = append(server.messageQueues[clientID], message)
	return nil
}

func RegisterWebSocketRoutes(e *echo.Echo, wsServer *WebSocketServer) {
	e.GET("/ws/:id", wsServer.handleWebSocket)
}
