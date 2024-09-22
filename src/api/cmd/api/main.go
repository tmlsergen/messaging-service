package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	_ "github.com/joho/godotenv/autoload"
	slogfiber "github.com/samber/slog-fiber"
	_ "github.com/tmlsergen/messaging-service-api/docs"
	app2 "github.com/tmlsergen/messaging-service-api/internal/app"
	"github.com/tmlsergen/messaging-service-api/internal/message"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:80
// @BasePath /
func main() {
	app := fiber.New(
		fiber.Config{
			ErrorHandler: app2.ErrorResponse,
		},
	)
	app.Use(requestid.New())

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	app.Use(slogfiber.NewWithConfig(logger, slogfiber.Config{
		DefaultLevel:       slog.LevelInfo,
		ClientErrorLevel:   slog.LevelWarn,
		ServerErrorLevel:   slog.LevelError,
		WithRequestBody:    true,
		WithResponseBody:   true,
		WithRequestHeader:  true,
		WithResponseHeader: true,
	}))

	// init db
	db := app2.InitDB()

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	// init redis
	rds := app2.InitRedis()

	route := app.Group("/api/v1")
	{
		//message routes
		messageRepository := message.NewRepository(db)
		messageService := message.NewService(messageRepository, rds)
		messageHandler := message.NewHandler(messageService)

		route.Get("/messages", messageHandler.GetMessages)
		route.Post("/messages/cron", messageHandler.HandleCronAction)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello, World!")
	})

	app.Get("/health-check", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Listen(":" + os.Getenv("APP_PORT"))
}
