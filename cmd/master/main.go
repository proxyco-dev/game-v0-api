package main

import (
	"context"
	_ "game-v0-api/docs"
	"game-v0-api/pkg/config"
	"game-v0-api/pkg/database"
	models "game-v0-api/pkg/models"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uptrace/bun"
)

type HandlerV1 struct {
	db *bun.DB
}

// @Summary Create a new room
// @Description Create a new room
// @Accept json
// @Produce json
// @Success 200 {object} models.Room
// @Router /api/v1/rooms [post]
func (this HandlerV1) CreateRoomV1(c *gin.Context) {
	dto := models.RoomDTO{MaxPlayers: 2, Private: false}

	if err := c.ShouldBind(&dto); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	slog.Info("DTO is {}", dto)

	// spin up game server => Address
	address := ""

	// create Room model => Generate a Code for the room
	code := "abcd"

	// save Room model in database => Fill in all fields that remain
	room := models.Room{Code: code, Address: address, Name: dto.Name, MaxPlayers: dto.MaxPlayers, Private: dto.Private}

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
// @Success 200 {array} models.Room
// @Router /api/v1/rooms [get]
func (this HandlerV1) ListRoomsV1(c *gin.Context) {
	rooms := []models.Room{}

	err := this.db.NewSelect().Model(&rooms).Where("? = ?", bun.Ident("private"), false).Scan(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, rooms)
}

// @Summary Join a room
// @Description Join a room
// @Accept json
// @Produce json
// @Param code path string true "Room code"
// @Success 200 {object} models.Room
// @Router /api/v1/rooms/join/{code} [get]
func (this HandlerV1) JoinRoomV1(c *gin.Context) {
	code := c.Param("code")
	room := models.Room{}

	err := this.db.NewSelect().Model(&room).Where("? = ?", bun.Ident("code"), code).Limit(1).Scan(c)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, room)
}

type RoomsService interface {
	CreateRoom(c context.Context, room models.RoomDTO) (models.Room, error)
}

type roomsService struct {
	// gameService GameService
	// roomsRepository RoomsRepository
}

func NewRoomService() RoomsService {
	return roomsService{}
}

func (this roomsService) CreateRoom(c context.Context, room models.RoomDTO) (models.Room, error) {
	return models.Room{Name: room.Name, MaxPlayers: room.MaxPlayers, Private: room.Private}, nil
}

// @title Game V0 API
// @version 1.0
// @description Game V0 API
// @host localhost:8080
// @BasePath /
// @contact.name Nika Shamiladze
// @contact.email fbshamiladze@gmail.com
func main() {
	config := config.LoadConfig()
	slog.Info("The config is: %v", config)

	db := database.New(config)

	router := gin.Default()

	handler := HandlerV1{db}

	apiV1 := router.Group("/api/v1")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	apiV1.GET("/rooms", handler.ListRoomsV1)
	apiV1.POST("/rooms", handler.CreateRoomV1)
	apiV1.GET("/rooms/join/:code", handler.JoinRoomV1)

	router.Run()
}
