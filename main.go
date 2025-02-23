package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
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

type Enemy struct {
	ID     string
	X, Y   float64
	Health int
	VX, VY float64
}

type GameState struct {
	Players map[string]*Player
	Bullets map[string]*Bullet
	Enemies map[string]*Enemy
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
	Enemies: make(map[string]*Enemy),
}

const serverSpeed = 150.0
const enemySpawnInterval = 6 * time.Second
const enemySpeed = 80.0

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
		go spawnEnemies()
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
					log.Println("Bullet spawned:", bulletID, "at x:", player.X, "y:", player.Y, "dirX:", dirX, "dirY:", dirY)
				}
			}
		}
	}
}

func spawnEnemies() {
	ticker := time.NewTicker(enemySpawnInterval)
	defer ticker.Stop()

	for range ticker.C {
		game.Mutex.Lock()
		enemyID := fmt.Sprintf("enemy%d", time.Now().UnixNano())
		x := float64(rand.Intn(640))
		y := float64(rand.Intn(480))
		if rand.Float64() < 0.5 {
			if rand.Float64() < 0.5 {
				x = -50
			} else {
				x = 690
			}
		} else {
			if rand.Float64() < 0.5 {
				y = -50
			} else {
				y = 530
			}
		}
		game.Enemies[enemyID] = &Enemy{
			ID:     enemyID,
			X:      x,
			Y:      y,
			Health: 30,
			VX:     0,
			VY:     0,
		}
		log.Println("Spawned enemy:", enemyID, "at x:", x, "y:", y)
		game.Mutex.Unlock()
	}
}

func broadcastGameState() {
	tickRate := 60
	ticker := time.NewTicker(time.Second / time.Duration(tickRate))
	defer ticker.Stop()

	for range ticker.C {
		game.Mutex.Lock()

		now := time.Now()
		deltaTime := time.Second.Seconds() / float64(tickRate)

		for _, p := range game.Players {
			pLastUpdate := p.LastUpdate
			if now.Sub(pLastUpdate).Seconds() > 0 {
				vx, vy := p.VelocityX, p.VelocityY
				if vx != 0 || vy != 0 {
					mag := math.Sqrt(vx*vx + vy*vy)
					vx /= mag
					vy /= mag
				}
				p.X += vx * serverSpeed * deltaTime
				p.Y += vy * serverSpeed * deltaTime
				p.LastUpdate = now
			}
		}

		for _, e := range game.Enemies {
			nearestPlayer := findNearestPlayer(e.X, e.Y)
			if nearestPlayer != nil {
				dx := nearestPlayer.X - e.X
				dy := nearestPlayer.Y - e.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist > 5 {
					e.VX = (dx / dist) * enemySpeed
					e.VY = (dy / dist) * enemySpeed
				} else {
					e.VX = 0
					e.VY = 0
				}
				e.X += e.VX * deltaTime
				e.Y += e.VY * deltaTime
			}
		}

		for id, b := range game.Bullets {
			prevX, prevY := b.X, b.Y
			b.X += b.VX
			b.Y += b.VY

			for eid, e := range game.Enemies {
				if lineIntersectsCircle(prevX, prevY, b.X, b.Y, e.X, e.Y, 10) {
					e.Health -= 10
					log.Println("Bullet", id, "hit enemy", eid, "new health:", e.Health)
					delete(game.Bullets, id)
					if e.Health <= 0 {
						delete(game.Enemies, eid)
						log.Println("Enemy killed:", eid)
					}
					break
				}
			}
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
			if b.X < -50 || b.X > 690 || b.Y < -50 || b.Y > 530 {
				delete(game.Bullets, id)
			}
		}

		state := map[string]interface{}{
			"players": make(map[string]map[string]interface{}),
			"bullets": make(map[string]map[string]interface{}),
			"enemies": make(map[string]map[string]interface{}),
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
		for id, e := range game.Enemies {
			state["enemies"].(map[string]map[string]interface{})[id] = map[string]interface{}{
				"x":      e.X,
				"y":      e.Y,
				"health": e.Health,
			}
		}
		log.Println("Broadcasting state - Players:", len(game.Players), "Enemies:", len(game.Enemies))
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

func lineIntersectsCircle(x1, y1, x2, y2, cx, cy, radius float64) bool {
	dx := x2 - x1
	dy := y2 - y1
	fx := x1 - cx
	fy := y1 - cy

	a := dx*dx + dy*dy
	b := 2 * (fx*dx + fy*dy)
	c := fx*fx + fy*fy - radius*radius
	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return false
	}

	t := (-b - math.Sqrt(discriminant)) / (2 * a)
	if t >= 0 && t <= 1 {
		return true
	}
	return false
}

func findNearestPlayer(ex, ey float64) *Player {
	var nearest *Player
	minDist := math.MaxFloat64
	for _, p := range game.Players {
		dist := distance(ex, ey, p.X, p.Y)
		if dist < minDist {
			minDist = dist
			nearest = p
		}
	}
	return nearest
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
