package main

import (
	"game-v0-api/api/handlers"
	"game-v0-api/database"
	"game-v0-api/pkg/common"
	"game-v0-api/pkg/redis"
	roomRepository "game-v0-api/pkg/room"
	userRepository "game-v0-api/pkg/user"
	"log"
	"os"

	_ "game-v0-api/docs"

	"game-v0-api/api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
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

	godotenv.Load()

	database.ConnectDB()

	app.Use(cors.New())
	bundle := i18n.NewBundle(language.Georgian)

	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	bundle.MustLoadMessageFile("./lang/active.ka.toml")
	bundle.MustLoadMessageFile("./lang/active.en.toml")

	app.Use(middleware.I18nMiddleware(bundle))

	redis.InitRedis(os.Getenv("REDIS_URL"), os.Getenv("REDIS_PASSWORD"))

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	api := app.Group("/api")

	// User Routes
	userRepo := userRepository.NewUserRepository(database.DB)
	userHandler := handlers.NewUserHandler(userRepo)

	api.Get("/user/me", common.AuthMiddleware, userHandler.GetMe)

	api.Post("/user/sign-in", userHandler.SignIn)
	api.Post("/user/sign-up", userHandler.SignUp)

	// Room Routes
	roomRepo := roomRepository.NewRoomRepository(database.DB)

	wsHandler := handlers.NewWebSocketHandler()
	go wsHandler.GetManager().Run()

	roomHandler := handlers.NewRoomHandler(roomRepo, bundle, wsHandler)

	api.Get("/room", common.AuthMiddleware, roomHandler.GetRooms)
	api.Get("/room/:id", common.AuthMiddleware, roomHandler.FindOne)

	api.Post("/room", common.AuthMiddleware, roomHandler.CreateRoom)
	api.Post("/room/join", common.AuthMiddleware, roomHandler.JoinRoom)
	api.Post("/room/leave", common.AuthMiddleware, roomHandler.LeaveRoom)

	api.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			roomID := c.Query("roomId")
			if roomID == "" {
				return fiber.ErrBadRequest
			}
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	api.Get("/ws", websocket.New(wsHandler.HandleWebSocket))

	log.Fatal(app.Listen(":8080"))
}
