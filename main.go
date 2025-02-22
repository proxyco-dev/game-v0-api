package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

type Player struct {
	ID         string
	X, Y       float64
	Health     int
	Conn       *websocket.Conn
	Mutex      sync.Mutex
	VelocityX  float64
	VelocityY  float64
	LastUpdate time.Time
}

type Bullet struct {
	ID      string
	X, Y    float64
	VX, VY  float64
	OwnerID string
}

type GameState struct {
	Players map[string]*Player
	Bullets map[string]*Bullet
	Mutex   sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var game = GameState{
	Players: make(map[string]*Player),
	Bullets: make(map[string]*Bullet),
}

const serverSpeed = 150.0

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWebSocket)

	handler := cors.Default().Handler(mux)

	port := ":8080"
	log.Printf("Server starting on %s...\n", port)
	log.Printf("WebSocket endpoint available at ws://localhost%s/ws\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	playerID := fmt.Sprintf("player%d", len(game.Players)+1)
	player := &Player{
		ID:         playerID,
		X:          100,
		Y:          100,
		Health:     100,
		Conn:       conn,
		LastUpdate: time.Now(),
		VelocityX:  0,
		VelocityY:  0,
	}
	game.Mutex.Lock()
	game.Players[playerID] = player
	log.Println("Added player:", playerID, "Total players:", len(game.Players))
	game.Mutex.Unlock()

	conn.WriteJSON(map[string]interface{}{
		"type": "init",
		"id":   playerID,
	})
	log.Println("Sent init to:", playerID)

	if len(game.Players) == 1 {
		go broadcastGameState()
	}

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Disconnect:", playerID, err)
			game.Mutex.Lock()
			delete(game.Players, playerID)
			log.Println("Removed player:", playerID, "Remaining:", len(game.Players))
			game.Mutex.Unlock()
			return
		}
		handleInput(player, msg)
	}
}

func handleInput(player *Player, msg map[string]interface{}) {
	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	if action, ok := msg["action"].(string); ok {
		switch action {
		case "move":
			if vx, ok := msg["velocityX"].(float64); ok {
				player.VelocityX = vx
			}
			if vy, ok := msg["velocityY"].(float64); ok {
				player.VelocityY = vy
			}
		case "shoot":
			if dirX, ok := msg["dirX"].(float64); ok {
				if dirY, ok := msg["dirY"].(float64); ok {
					bulletID := fmt.Sprintf("bullet%d", len(game.Bullets)+1)
					game.Bullets[bulletID] = &Bullet{
						ID:      bulletID,
						X:       player.X,
						Y:       player.Y,
						VX:      dirX * 5,
						VY:      dirY * 5,
						OwnerID: player.ID,
					}
				}
			}
		}
	}
}
func broadcastGameState() {
	tickRate := 60
	ticker := time.NewTicker(time.Second / time.Duration(tickRate))
	defer ticker.Stop()

	for range ticker.C {
		game.Mutex.Lock()

		now := time.Now()
		for _, p := range game.Players {
			deltaTime := now.Sub(p.LastUpdate).Seconds()
			if deltaTime > 0 {
				vx, vy := p.VelocityX, p.VelocityY
				if vx != 0 || vy != 0 {
					mag := math.Sqrt(vx*vx + vy*vy)
					vx /= mag
					vy /= mag
					oldX, oldY := p.X, p.Y
					p.X += vx * serverSpeed * deltaTime
					p.Y += vy * serverSpeed * deltaTime
					log.Printf("Player %s moved from x:%.2f, y:%.2f to x:%.2f, y:%.2f", p.ID, oldX, oldY, p.X, p.Y)
				}
				p.LastUpdate = now
			}
		}

		for id, b := range game.Bullets {
			b.X += b.VX
			b.Y += b.VY
			for _, p := range game.Players {
				if p.ID != b.OwnerID && distance(b.X, b.Y, p.X, p.Y) < 10 {
					p.Health -= 10
					delete(game.Bullets, id)
					if p.Health <= 0 {
						delete(game.Players, p.ID)
						p.Conn.Close()
					}
					break
				}
			}
			if b.X < 0 || b.X > 640 || b.Y < 0 || b.Y > 480 {
				delete(game.Bullets, id)
			}
		}

		state := map[string]interface{}{
			"players": make(map[string]map[string]interface{}),
			"bullets": make(map[string]map[string]interface{}),
		}
		for id, p := range game.Players {
			state["players"].(map[string]map[string]interface{})[id] = map[string]interface{}{
				"x":      p.X,
				"y":      p.Y,
				"health": p.Health,
			}
		}
		for id, b := range game.Bullets {
			state["bullets"].(map[string]map[string]interface{})[id] = map[string]interface{}{
				"x": b.X,
				"y": b.Y,
			}
		}
		log.Println("Broadcasting state - Players:", len(game.Players))
		game.Mutex.Unlock()

		for _, p := range game.Players {
			p.Mutex.Lock()
			err := p.Conn.WriteJSON(state)
			p.Mutex.Unlock()
			if err != nil {
				log.Println("Error broadcasting to", p.ID, ":", err)
			}
		}
	}
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
