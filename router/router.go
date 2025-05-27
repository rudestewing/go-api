package router

import (
	"go-api/container"
	"go-api/internal/handler"
	"go-api/internal/middleware"
	"go-api/shared/response"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, container *container.Container) {
	v1 := app.Group("/api/v1")
	v1.Get("/", func(c *fiber.Ctx) error {
		return response.SuccessWithMessage(c, "hello from v1", "API is running")
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
