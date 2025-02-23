package handlers

import (
	"game-v0-api/api/presenter"
	repository "game-v0-api/pkg/room"

	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type RoomHandler struct {
	roomRepo repository.RoomRepository
	bundle   *i18n.Bundle
}

func NewRoomHandler(roomRepo repository.RoomRepository, bundle *i18n.Bundle) *RoomHandler {
	return &RoomHandler{
		roomRepo: roomRepo,
		bundle:   bundle,
	}
}

func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Room created successfully",
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
