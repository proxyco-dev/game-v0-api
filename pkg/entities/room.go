package entities

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Users      []*User    `json:"users" gorm:"many2many:room_users;"`
	Title      string     `json:"title"`
	IsActive   bool       `json:"isActive" gorm:"default:true"`
	MaxPlayers int        `json:"maxPlayers" gorm:"default:2"`
	IsDeleted  bool       `json:"isDeleted" gorm:"default:false"`
	CreatedAt  *time.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt" gorm:"column:createdAt"`
}

func (Room) TableName() string {
	return "rooms-v0"
}
