package main

import (
	"context"
	matchmaker "game-v0-api/cmd"
	_ "game-v0-api/docs"
	"log/slog"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uptrace/bun"
)

type GameServer interface {
	Spawn(code string, maxPlayers int) (string, int, int, error)
}

type gameServerDocker struct {
}

func (this gameServerDocker) Spawn(code string, maxPlayers int) (address string, query int, game int, err error) {

	listener, err := net.Listen("tcp", ":0")

	if err != nil {

		return

	}

	query = listener.Addr().(*net.TCPAddr).Port

	listener.Close()

	listener, err = net.Listen("tcp", ":0")

	if err != nil {

		return

	}

	game = listener.Addr().(*net.TCPAddr).Port

	listener.Close()

	address = "0.0.0.0"

	// TODO: run the docker command

	return

}

func NewGameServer(cfg matchmaker.Config) GameServer {

	return gameServerDocker{}

}

type HandlerV1 struct {
	db         *bun.DB
	gameServer GameServer
}

// @Summary Create a new room
// @Description Create a new room
// @Accept json
// @Produce json
// @Success 200 {object} matchmaker.Room
// @Router /api/v1/rooms [post]
func (this HandlerV1) CreateRoom(c *gin.Context) {
	dto := matchmaker.RoomDTO{MaxPlayers: 2, Private: false}

	if err := c.ShouldBind(&dto); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	slog.Info("DTO is {}", dto)

	code := "abcd"

	address, queryPort, gamePort, err := this.gameServer.Spawn(code, dto.MaxPlayers)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	room := matchmaker.Room{
		Code:       code,
		Address:    address,
		QueryPort:  queryPort,
		GamePort:   gamePort,
		Name:       dto.Name,
		MaxPlayers: dto.MaxPlayers,
		Private:    dto.Private,
	}

	_, err = this.db.NewInsert().Model(&room).Exec(c)

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
// @Success 200 {array} matchmaker.Room
// @Router /api/v1/rooms [get]
func (this HandlerV1) ListRooms(c *gin.Context) {
	rooms := []matchmaker.Room{}

	err := this.db.NewSelect().Model(&rooms).Where("? = ?", bun.Ident("private"), false).Scan(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, rooms)
}

// @Summary Get a room
// @Description Get a room
// @Accept json
// @Produce json
// @Param code path string true "Room code"
// @Success 200 {object} matchmaker.Room
// @Router /api/v1/rooms/{code} [get]
func (this HandlerV1) GetRoom(c *gin.Context) {
	code := c.Param("code")
	room := matchmaker.Room{}

	err := this.db.NewSelect().Model(&room).Where("? = ?", bun.Ident("code"), code).Limit(1).Scan(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, room)
}

type RoomsService interface {
	CreateRoom(c context.Context, room matchmaker.RoomDTO) (matchmaker.Room, error)
}

type roomsService struct {
	// gameService GameService
	// roomsRepository RoomsRepository
}

func NewRoomService() RoomsService {
	return roomsService{}
}

func (this roomsService) CreateRoom(c context.Context, room matchmaker.RoomDTO) (matchmaker.Room, error) {
	return matchmaker.Room{Name: room.Name, MaxPlayers: room.MaxPlayers, Private: room.Private}, nil
}

// @title Game V0 API
// @version 1.0
// @description Game V0 API
// @host localhost:8080
// @BasePath /
// @contact.name Nika Shamiladze
// @contact.email fbshamiladze@gmail.com
func main() {
	config := matchmaker.LoadConfig()
	slog.Info("The config is: %v", config)

	db := matchmaker.NewDB(config)
	gameServer := NewGameServer(config)

	router := gin.Default()

	handler := HandlerV1{db, gameServer}

	apiV1 := router.Group("/api/v1")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	apiV1.GET("/rooms", handler.ListRooms)
	apiV1.POST("/rooms", handler.CreateRoom)
	apiV1.GET("/rooms/:code", handler.GetRoom)

	router.Run()
}
