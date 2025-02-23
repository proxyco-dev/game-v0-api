package main

import (
	"game-v0-api/api/handlers"
	"game-v0-api/database"
	roomRepository "game-v0-api/pkg/room"
	userRepository "game-v0-api/pkg/user"
	"log"

	_ "game-v0-api/docs"

	"game-v0-api/api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"golang.org/x/text/language"
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
	bundle := i18n.NewBundle(language.Georgian)

	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	bundle.MustLoadMessageFile("./lang/active.ka.toml")
	bundle.MustLoadMessageFile("./lang/active.en.toml")

	app.Use(middleware.I18nMiddleware(bundle))

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// User Routes
	userRepo := userRepository.NewUserRepository(database.DB)
	userHandler := handlers.NewUserHandler(userRepo)

	app.Get("/user/me", userHandler.GetMe)
	app.Post("/user/sign-in", userHandler.SignIn)
	app.Post("/user/sign-up", userHandler.SignUp)

	roomRepo := roomRepository.NewRoomRepository(database.DB)
	roomHandler := handlers.NewRoomHandler(roomRepo, bundle)

	app.Post("/room", roomHandler.CreateRoom)
	app.Get("/room", roomHandler.GetRooms)

	log.Fatal(app.Listen(":8080"))
}
