package websocket

import (
	"log"
	"sync"
)

// Hub хранит всех активных клиентов и управляет рассылкой сообщений
type Hub struct {
	clients    map[*Client]bool      // активные клиенты
	broadcast  chan []byte           // канал для входящих сообщений
	register   chan *Client          // канал для регистрации клиентов
	unregister chan *Client          // канал для отключения клиентов
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("✅ Клиент подключился: %s", client.ID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("❌ Клиент отключился: %s", client.ID)

		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}