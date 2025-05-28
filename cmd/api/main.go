package main

import (
	"go-api/config"
	"go-api/container"
	"go-api/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func createFiberApp() *fiber.App {
	cfg := config.Get()

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	return app
}

func main() {
	// Initialize config first
	config.InitConfig()

	app := createFiberApp()

	container, err := container.NewContainer()
	if err != nil {
		log.Fatal("Failed to create container:", err)
	}

	router.RegisterRoutes(app, container)

	cfg := config.Get()
	port := ":" + cfg.AppPort
	log.Printf("🚀 Server starting on port %s...", cfg.AppPort)

	if err := app.Listen(port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
