package main

import (
	"context"
	"go-api/config"
	"go-api/container"
	"go-api/router"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func createFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
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

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")

		// Shutdown server
		if err := app.Shutdown(); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}

		// Cleanup container
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := container.Close(ctx); err != nil {
			log.Printf("Error during container cleanup: %v", err)
		}

		os.Exit(0)
	}()

	router.RegisterRoutes(app, container)

	cfg := config.Get()
	port := ":" + cfg.AppPort
	log.Printf("Server starting on port %s...", cfg.AppPort)

	if err := app.Listen(port); err != nil {
		log.Printf("Server stopped: %v", err)
	}
}
