package handlers

import (
	"game-v0-api/pkg/websocket"
	"log"
	"time"

	ws "github.com/gofiber/websocket/v2"
)

type WebSocketHandler struct {
	manager *websocket.Manager
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		manager: websocket.NewManager(),
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *ws.Conn) {
	const (
		writeWait  = 10 * time.Second
		pongWait   = 60 * time.Second
		pingPeriod = (pongWait * 9) / 10
	)

	roomID := c.Query("roomId")
	if roomID == "" {
		log.Println("No room ID provided")
		return
	}

	client := &websocket.Client{
		Conn:   c,
		RoomID: roomID,
	}

	log.Printf("New client connected to room: %s", roomID)

	h.manager.Register <- client

	defer func() {
		h.manager.Unregister <- client
		if err := c.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := c.WriteControl(ws.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
					log.Printf("ping error: %v", err)
					return
				}
			}
		}
	}()

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if ws.IsCloseError(err,
				ws.CloseGoingAway,
				ws.CloseAbnormalClosure,
				ws.CloseNoStatusReceived) {
				log.Printf("Connection closed: %v", err)
			} else {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		if messageType == ws.TextMessage {
			log.Printf("Received message: %s", message)
		}
	}
}

func (h *WebSocketHandler) GetManager() *websocket.Manager {
	return h.manager
}
