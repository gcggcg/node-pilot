package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"node-pilot/internal/model"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	TaskID uint64
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

func (h *Hub) Run() {
	log.Printf("[WS-HUB] Hub Run loop started")
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			log.Printf("[WS-HUB] Client registered, total clients: %d", len(h.clients))
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Printf("[WS-HUB] Client unregistered, total clients: %d", len(h.clients))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			log.Printf("[WS-HUB] Broadcast received, sending to %d clients", len(h.clients))
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Broadcast(msg *model.WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("broadcast marshal error: %v", err)
		return
	}

	h.mu.RLock()
	clientCount := len(h.clients)
	h.mu.RUnlock()

	log.Printf("[WS-BROADCAST] type=%s, task_id=%d, server_id=%d, client_count=%d",
		msg.Type, msg.TaskID, msg.ServerID, clientCount)

	if clientCount == 0 {
		log.Printf("[WS-BROADCAST] No clients connected, message dropped")
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- data:
			log.Printf("[WS-BROADCAST] Sent to client")
		default:
			close(client.Send)
			delete(h.clients, client)
			log.Printf("[WS-BROADCAST] Client send buffer full, removed")
		}
	}
}

func (h *Hub) BroadcastToTask(msg *model.WSMessage, taskID uint64) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("broadcast marshal error: %v", err)
		return
	}

	h.mu.RLock()

	targetClients := make([]*Client, 0)
	for client := range h.clients {
		if client.TaskID == taskID {
			targetClients = append(targetClients, client)
		}
	}
	clientCount := len(targetClients)
	h.mu.RUnlock()

	log.Printf("[WS-BROADCAST-TO-TASK] type=%s, task_id=%d, target_client_count=%d",
		msg.Type, msg.TaskID, clientCount)

	if clientCount == 0 {
		log.Printf("[WS-BROADCAST-TO-TASK] No clients watching task %d, message dropped", taskID)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for _, client := range targetClients {
		select {
		case client.Send <- data:
			log.Printf("[WS-BROADCAST-TO-TASK] Sent to client (task=%d)", taskID)
		default:
			close(client.Send)
			delete(h.clients, client)
			log.Printf("[WS-BROADCAST-TO-TASK] Client send buffer full, removed")
		}
	}
}

func (c *Client) ReadPump(taskID string) {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
