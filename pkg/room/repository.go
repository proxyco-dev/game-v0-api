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
	FindByIdWithUsers(id string) (*entities.Room, error)
	Update(room *entities.Room) error
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
	err := r.db.Table(roomsTableName).Preload("Users").Where("is_active = ? AND is_deleted = ? AND status = ?", true, false, "waiting").Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepository) FindById(id string) (*entities.Room, error) {
	var room entities.Room
	err := r.db.Table(roomsTableName).Preload("Users").Where("id = ? AND is_deleted = ? AND is_active = ?", id, false, true).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) FindByIdWithUsers(id string) (*entities.Room, error) {
	var room entities.Room
	err := r.db.Table(roomsTableName).Where("id = ? AND is_deleted = ? AND is_active = ?", id, false, true).Preload("Users").First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) Update(room *entities.Room) error {
	return r.db.Table(roomsTableName).Save(room).Error
}
