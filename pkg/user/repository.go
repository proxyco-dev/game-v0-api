package repository

import (
	models "game-v0-api/pkg/entities"

	"gorm.io/gorm"
)

const (
	usersTableName = "users-v0"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Table(usersTableName).Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Table(usersTableName).First(&user, "\"email\" = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
