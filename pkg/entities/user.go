package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:createdAt"`
}

func (User) TableName() string {
	return "users-v0"
}
