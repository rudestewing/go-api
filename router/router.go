package router

import (
	"go-api/app/handler"
	"go-api/container"
	"go-api/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, container *container.Container) {
	router := app.Group("/api/v1")
	router.Use(middleware.TimeoutMiddleware(30*time.Second, "Operation timed out"))
	router.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"data":    nil,
			"message": "API is running",
		})
	})

	// register handlers here...
	handler.RegisterAuthHandler(router, container)
}
