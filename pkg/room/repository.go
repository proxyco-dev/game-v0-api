package repository

import (
	"game-v0-api/pkg/entities"

	"gorm.io/gorm"
)

const (
	roomsTableName = "rooms-v0"
)

type RoomRepository interface {
	Create(room *entities.Room) error
	FindAll() ([]entities.Room, error)
	FindById(id string) (*entities.Room, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(room *entities.Room) error {
	return r.db.Table(roomsTableName).Create(room).Error
}

func (r *roomRepository) FindAll() ([]entities.Room, error) {
	var rooms []entities.Room
	err := r.db.Table(roomsTableName).Find(&rooms).Where("is_active = ? AND is_deleted = ?", true, false).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepository) FindById(id string) (*entities.Room, error) {
	var room entities.Room
	err := r.db.Table(roomsTableName).Where("id = ? AND is_deleted = ? AND is_active = ?", id, false, true).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}
