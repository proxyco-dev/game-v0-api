package handlers

import (
	"context"
	"encoding/json"
	"game-v0-api/api/presenter"
	entities "game-v0-api/pkg/entities"
	"game-v0-api/pkg/redis"
	repository "game-v0-api/pkg/room"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type RoomHandler struct {
	roomRepo         repository.RoomRepository
	bundle           *i18n.Bundle
	validator        *validator.Validate
	websocketHandler *WebSocketHandler
}

func NewRoomHandler(roomRepo repository.RoomRepository, bundle *i18n.Bundle, websocketHandler *WebSocketHandler) *RoomHandler {
	return &RoomHandler{
		roomRepo:         roomRepo,
		bundle:           bundle,
		validator:        validator.New(),
		websocketHandler: websocketHandler,
	}
}

// @Summary Create a room
// @Description Create a room
// @Tags Room
// @Accept json
// @Produce json
// @Param room body presenter.RoomRequest true "Room"
// @Success 200 {object} presenter.RoomResponse
// @Router /api/room [post]
func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	var request presenter.RoomRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(presenter.ErrorResponse{
			Error: err.Error(),
		})
	}

	if err := h.validator.Struct(request); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(presenter.ErrorResponse{Error: validationErrors.Error()})
	}

	room := &entities.Room{
		Title:      request.Title,
		MaxPlayers: request.MaxPlayers,
	}

	if err := h.roomRepo.Create(room); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse{Error: err.Error()})
	}

	redisClient := redis.GetClient()
	ctx := context.Background()

	roomJson, err := json.Marshal(room)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse{Error: err.Error()})
	}

	err = redisClient.Set(ctx, "rooms:"+room.ID.String(), roomJson, 0).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse{Error: err.Error()})
	}

	localizer := c.Locals("localizer").(*i18n.Localizer)

	message, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "CreatedSuccessfully",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}

// @Summary Get all rooms
// @Description Get all rooms
// @Tags Room
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/room [get]
func (h *RoomHandler) GetRooms(c *fiber.Ctx) error {
	localizer := c.Locals("localizer").(*i18n.Localizer)

	rooms, err := h.roomRepo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse{
			Error: err.Error(),
		})
	}

	message, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "FetchedSuccessfully",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"data":    rooms,
	})
}

// @Summary Join a room
// @Description Join a room
// @Tags Room
// @Accept json
// @Produce json
// @Param room body presenter.JoinRoomRequest true "Room"
// @Success 200 {object} presenter.RoomResponse
// @Router /api/room/join [post]
func (h *RoomHandler) JoinRoom(c *fiber.Ctx) error {
	var request presenter.JoinRoomRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(presenter.ErrorResponse{Error: err.Error()})
	}

	if err := h.validator.Struct(request); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(presenter.ErrorResponse{Error: validationErrors.Error()})
	}

	redisClient := redis.GetClient()
	ctx := context.Background()

	var room *entities.Room
	redisRoom, err := redisClient.Get(ctx, "rooms:"+request.Id).Result()
	if err != nil {
		room, err = h.roomRepo.FindById(request.Id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(presenter.ErrorResponse{Error: err.Error()})
		}
	} else {
		if err := json.Unmarshal([]byte(redisRoom), &room); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse{Error: err.Error()})
		}
	}

	id := c.Locals("user").(jwt.MapClaims)["id"].(string)

	room.Users = append(room.Users, &entities.User{ID: uuid.MustParse(id)})

	if err := h.roomRepo.Update(room); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse{Error: err.Error()})
	}

	redisClient.Del(ctx, "rooms:"+room.ID.String())

	localizer := c.Locals("localizer").(*i18n.Localizer)

	message, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "JoinedSuccessfully",
	})

	h.websocketHandler.GetManager().EmitToRoom(room.ID.String(), "USER_JOINED", fiber.Map{
		"message": message,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"room":    room,
	})
}

// @Summary Find one room
// @Description Find one room
// @Tags Room
// @Accept json
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} presenter.RoomResponse
// @Router /api/room/{id} [get]
func (h *RoomHandler) FindOne(c *fiber.Ctx) error {
	id := c.Params("id")
	redisClient := redis.GetClient()
	ctx := context.Background()

	var room *entities.Room
	redisRoom, err := redisClient.Get(ctx, "rooms:"+id).Result()
	if err != nil {
		room, err = h.roomRepo.FindByIdWithUsers(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(presenter.ErrorResponse{Error: err.Error()})
		}

		roomJson, err := json.Marshal(room)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse{Error: err.Error()})
		}

		redisClient.Set(ctx, "rooms:"+id, roomJson, 0)
	} else {
		if err := json.Unmarshal([]byte(redisRoom), &room); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse{Error: err.Error()})
		}
	}

	localizer := c.Locals("localizer").(*i18n.Localizer)

	message, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "FetchedSuccessfully",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"room":    room,
	})
}
