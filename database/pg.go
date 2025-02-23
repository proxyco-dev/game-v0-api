package database

import (
	"fmt"
	models "game-v0-api/pkg/entities"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	return db, nil
}

func MigrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.User{}, &models.Room{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
