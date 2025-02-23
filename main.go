package main

import (
	"game-v0-api/api/handlers"
	"game-v0-api/database"
	roomRepository "game-v0-api/pkg/room"
	userRepository "game-v0-api/pkg/user"
	"log"

	_ "game-v0-api/docs"

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

	// User Routes
	userRepo := userRepository.NewUserRepository(database.DB)
	userHandler := handlers.NewUserHandler(userRepo)

	app.Get("/user/me", userHandler.GetMe)
	app.Post("/user/sign-in", userHandler.SignIn)
	app.Post("/user/sign-up", userHandler.SignUp)

	roomRepo := roomRepository.NewRoomRepository(database.DB)
	roomHandler := handlers.NewRoomHandler(roomRepo)

	app.Post("/room", roomHandler.CreateRoom)
	app.Get("/room", roomHandler.GetRooms)

	log.Fatal(app.Listen(":8080"))
}
