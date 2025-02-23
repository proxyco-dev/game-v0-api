package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn   *websocket.Conn
	RoomID string
}

type Event struct {
	Type    string      `json:"type"`
	RoomID  string      `json:"roomId"`
	Payload interface{} `json:"payload"`
}

type Manager struct {
	clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	mutex      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.Register:
			m.mutex.Lock()
			m.clients[client] = true
			m.mutex.Unlock()

		case client := <-m.Unregister:
			m.mutex.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				client.Conn.Close()
			}
			m.mutex.Unlock()
		}
	}
}

func (m *Manager) EmitToRoom(roomID string, eventType string, payload interface{}) {
	event := Event{
		Type:    eventType,
		RoomID:  roomID,
		Payload: payload,
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	clientCount := 0
	for client := range m.clients {
		if client.RoomID == roomID {
			clientCount++
			if err := client.Conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Printf("Error sending message to client: %v", err)
				go func(c *Client) {
					m.Unregister <- c
				}(client)
			} else {
				log.Printf("Successfully sent message to client in room %s", roomID)
			}
		}
	}
	log.Printf("Emitted message to %d clients in room %s", clientCount, roomID)
}
