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
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type RoomHandler struct {
	roomRepo  repository.RoomRepository
	bundle    *i18n.Bundle
	validator *validator.Validate
}

func NewRoomHandler(roomRepo repository.RoomRepository, bundle *i18n.Bundle) *RoomHandler {
	return &RoomHandler{
		roomRepo:  roomRepo,
		bundle:    bundle,
		validator: validator.New(),
	}
}

// @Summary Create a room
// @Description Create a room
// @Tags Room
// @Accept json
// @Produce json
// @Param room body presenter.RoomRequest true "Room"
// @Success 200 {object} presenter.RoomResponse
// @Router /room [post]
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
		Title: request.Title,
	}

	if err := h.roomRepo.Create(room); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse{Error: err.Error()})
	}

	redisClient := redis.GetClient()
	ctx := context.Background()

	redisClient.Set(ctx, "rooms:"+room.ID.String(), room, 0)

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
// @Router /room [get]
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
// @Router /room/join [post]
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

	user := c.Locals("user").(*entities.User)

	room.Users = append(room.Users, user)

	if err := h.roomRepo.Update(room); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse{Error: err.Error()})
	}

	redisClient.Set(ctx, "rooms:"+room.ID.String(), room, 0)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Room joined successfully",
		"room":    room,
	})
}
