package router

import (
	"go-api/container"
	"go-api/internal/handler"
	"go-api/internal/middleware"
	"go-api/shared/response"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, container *container.Container) {
	v1 := app.Group("/api/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return response.SuccessWithMessage(c, "hello from v1", "API is running")
	})

	// Apply timeout middleware
	v1.Use(middleware.TimeoutMiddlewareWithCustomMessage(30*time.Second, "User operation timed out"))

	authHandler := handler.NewAuthHandler(container.AuthService)
	// Auth routes (public) - no timeout middleware for login/register
	auth := v1.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)

	// Protected routes with timeout middleware
	userHandler := handler.NewUserHandler(container.UserService)
	protected := v1.Group("/user")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.Get("/profile", userHandler.GetProfile)
}
