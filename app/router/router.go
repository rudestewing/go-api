package router

import (
	"go-api/app/handler"
	"go-api/app/middleware"
	"go-api/app/shared/response"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler *handler.Handler) {
	v1 := app.Group("/api/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return response.Success(c, nil, "hello from v1", "API is running")
	})

	// Apply timeout middleware
	v1.Use(middleware.TimeoutMiddleware(30*time.Second, "User operation timed out"))

	// Auth routes (public) - no timeout middleware for login/register
	auth := v1.Group("/auth")
	auth.Post("/register", handler.AuthHandler.Register)
	auth.Post("/login", handler.AuthHandler.Login)

	// Protected routes with timeout middleware
	protected := v1.Group("/user")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.Get("/profile", handler.UserHandler.GetProfile)


	v1.Use(middleware.TimeoutMiddleware(30 * time.Second))
	v1.Get("asynchronous", handler.AsynchronousHandler.RunAsyncTask)
}
