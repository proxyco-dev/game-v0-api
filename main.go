package main

import (
	"context"
	_ "game-v0-api/docs"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uptrace/bun"
)

type RoomDTO struct {
	Name       string `json:"name" binding:"required"`
	MaxPlayers int    `json:"maxPlayers"`
	Private    bool   `json:"private"`
}

type Room struct {
	bun.BaseModel `bun:"table:rooms,alias:r"`

	ID         int64     `bun:"id,pk,autoincrement" json:"id"`
	Code       string    `bun:"code,type:varchar(6),notnull" json:"code"`
	Address    string    `bun:"address,type:varchar(128),notnull" json:"address"`
	Name       string    `bun:"name,type:varchar(128),notnull" json:"name"`
	MaxPlayers int       `bun:"maxPlayers,type:int,notnull,default:2" json:"maxPlayers"`
	Private    bool      `bun:"private,notnull,default:false"`
	CreatedAt  time.Time `bun:"createdAt,nullzero,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time `bun:"updatedAt,nullzero,notnull,default:current_timestamp" json:"updatedAt"`
}

type HandlerV1 struct {
	db *bun.DB
}

// @Summary Create a new room
// @Description Create a new room
// @Accept json
// @Produce json
// @Param room body RoomDTO true "Room details"
// @Success 200 {object} RoomDTO
// @Router /api/v1/rooms [post]
func (this HandlerV1) CreateRoomV1(c *gin.Context) {
	dto := RoomDTO{MaxPlayers: 2, Private: false}

	if err := c.ShouldBind(&dto); err != nil {

		c.AbortWithError(http.StatusBadRequest, err)

		return

	}

	slog.Info("DTO is {}", dto)

	room := Room{Name: dto.Name, MaxPlayers: dto.MaxPlayers, Private: dto.Private}

	// spin up game server => Address

	// create Room model => Generate a Code for the room

	// save Room model in database => Fill in all fields that remain

	_, err := this.db.NewInsert().Model(&room).Exec(c)

	if err != nil {

		c.AbortWithError(http.StatusBadRequest, err)

		return

	}

	c.JSON(http.StatusOK, room)
}

// @Summary List all rooms
// @Description List all rooms
// @Accept json
// @Produce json
// @Success 200 {array} RoomDTO
// @Router /api/v1/rooms [get]
func (this HandlerV1) ListRoomsV1(c *gin.Context) {
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
func (this HandlerV1) JoinRoomV1(c *gin.Context) {
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

func (this roomsService) CreateRoom(c context.Context, room RoomDTO) (Room, error) {
	return Room{Name: room.Name, MaxPlayers: room.MaxPlayers, Private: room.Private}, nil
}

// @title Game V0 API
// @version 1.0
// @description Game V0 API
// @host localhost:8080
// @BasePath /
// @contact.name Nika Shamiladze
// @contact.email fbshamiladze@gmail.com
func main() {
	router := gin.Default()

	handler := HandlerV1{db: nil}

	apiV1 := router.Group("/api/v1")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	apiV1.GET("/rooms", handler.ListRoomsV1)
	apiV1.POST("/rooms", handler.CreateRoomV1)
	apiV1.GET("/join/:code", handler.JoinRoomV1)

	router.Run()
}
