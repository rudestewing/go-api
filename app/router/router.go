package router

import (
	"go-api/app/handler"
	"go-api/app/middleware"
	"go-api/app/shared/response"
	"go-api/container"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, container *container.Container) {
	// Global middlewares
	app.Use(middleware.RecoverMiddleware())
	app.Use(middleware.ErrorHandler())
	
	// API routes
	router := app.Group("/api/v1")
	router.Use(middleware.TimeoutMiddleware(30*time.Second, "Operation timed out"))

	router.Get("/", func(c *fiber.Ctx) error {
		return response.Success(c, nil, "API is running", "Welcome to Go API v1")
	})

	// Health check endpoint
	router.Get("/health", func(c *fiber.Ctx) error {
		return response.Success(c, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "1.0.0",
		}, "Service is healthy")
	})

	// Register handlers here...
	handler.RegisterAuthHandler(router, container)
	
	// Handle 404 for API routes
	router.Use(middleware.NotFoundHandler())
}
