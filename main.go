package main

import (
	"context"
	_ "game-v0-api/docs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uptrace/bun"
)

type RoomDTO struct {
	Name string `json:"name" binding:"required"`
	MaxPlayers int	`json:"maxPlayers"`
	Description string `json:"description"`
	Private bool `json:"private"`
}

type Room struct {
	bun.BaseModel `bun:"table:rooms,alias:r"`

	ID int64 `bun:"id,pk,autoincrement" json:"id"`
	Address string `bun:"address,type:varchar(128),notnull" json:"address"`
	Name string `bun:"name,type:varchar(128),notnull" json:"name"`
	MaxPlayers int	`bun:"maxPlayers,type:int,notnull,default:2" json:"maxPlayers"`
	Description string `bun:"description,type:varchar(512),nullzero,notnull,default:''" json:"description"`
	Private bool `bun:"private,notnull,default:false" json:"private"`
	CreatedAt time.Time `bun:"createdAt,nullzero,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `bun:"updatedAt,nullzero,notnull,default:current_timestamp" json:"updatedAt"`
}

func NewRoom(r RoomDTO) Room {
	return Room{
		Name: r.Name,
		MaxPlayers: r.MaxPlayers,
		Description: r.Description,
		Private: r.Private,
	}
}

// @Summary List all rooms
// @Description List all rooms
// @Accept json
// @Produce json
// @Success 200 {array} RoomDTO
// @Router /api/v1/rooms [get]
func ListRoomsV1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// @Summary Create a new room
// @Description Create a new room
// @Accept json
// @Produce json
// @Param room body RoomDTO true "Room details"
// @Success 200 {object} RoomDTO
// @Router /api/v1/rooms [post]
func CreateRoomV1(c *gin.Context) {
	dto := RoomDTO{}

	if err := c.ShouldBind(&dto); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// @Summary Join a room
// @Description Join a room
// @Accept json
// @Produce json
// @Success 200 {object} RoomDTO
// @Router /api/v1/rooms/{ID}/join [post]
func JoinRoomV1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

type RoomsService interface {
	CreateRoom(c context.Context, room RoomDTO) (Room, error)
}

type roomsService struct {
	// gameService GameService
	// roomsRepository RoomsRepository
}

func NewRoomService() RoomsService {
	return roomsService{}
}

func (this roomsService)CreateRoom(c context.Context, room RoomDTO) (Room, error) {
 return Room{Name: room.Name, MaxPlayers: room.MaxPlayers, Description: room.Description, Private: room.Private}, nil
}

// @title Game V0 API
// @version 1.0
// @description Game V0 API
// @host localhost:8080
// @BasePath /
// @contact.name Nika Shamiladze
// @contact.email fbshamiladze@gmail.com
func main() {
	router := gin.Default();
	apiV1 := router.Group("/api/v1")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	apiV1.GET("/rooms", ListRoomsV1)
	apiV1.POST("/rooms", CreateRoomV1)
	apiV1.GET("/rooms/{ID}/join", JoinRoomV1)

	router.Run(":8080")
}
