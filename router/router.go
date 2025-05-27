package router

import (
	"go-api/container"
	"go-api/internal/handler"
	"go-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, container *container.Container) {
	v1 := app.Group("/api/v1")
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": "hello from v1",
		})
	})

	// Auth routes (public)
	auth := v1.Group("/auth")
	auth.Post("/login", container.AuthHandler.Login)
	auth.Post("/register", container.AuthHandler.Register)

	// Protected routes
	userHandler := handler.NewUserHandler()
	protected := v1.Group("/user")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.Get("/profile", userHandler.GetProfile)
}
