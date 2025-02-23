package handlers

import (
	"game-v0-api/api/presenter"
	repository "game-v0-api/pkg/room"

	"github.com/gofiber/fiber/v2"
)

type RoomHandler struct {
	roomRepo repository.RoomRepository
}

func NewRoomHandler(roomRepo repository.RoomRepository) *RoomHandler {
	return &RoomHandler{
		roomRepo: roomRepo,
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
	rooms, err := h.roomRepo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse{
			Error: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Rooms fetched successfully",
		"data":    rooms,
	})
}
