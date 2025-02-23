package main

import (
	"game-v0-api/database"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Game v0 API
// @version 1.0
// @description API for the survival game
// @contact.name Nika Shamiladze
// @contact.email fbshamiladze@gmail.com
// @host localhost:8080
// @BasePath /
func main() {
	app := fiber.New()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()

	app.Use(cors.New())

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":8080"))
}
