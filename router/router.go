package router

import (
	"go-api/app/handler"
	"go-api/container"
	"go-api/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, h *handler.Handler, container *container.Container) {
	// Create health handler directly from package
	healthHandler := handler.NewHealthHandler(container.DB)

	// Health check endpoints (no rate limiting for monitoring)
	app.Get("/health", healthHandler.HealthCheck)
	app.Get("/health/detailed", healthHandler.DetailedHealthCheck)
	app.Get("/health/ready", healthHandler.ReadinessCheck)
	app.Get("/health/live", healthHandler.LivenessCheck)

	router := app.Group("/api/v1")
	router.Use(middleware.TimeoutMiddleware(30*time.Second, "Operation timed out"))
	router.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"data":    nil,
			"message": "API is running",
		})
	})

	// Auth routes with specific rate limiting
	auth := router.Group("/auth")
	auth.Use(middleware.AuthRateLimitMiddleware()) // More restrictive rate limiting for auth
	auth.Post("/login", h.AuthHandler.Login)
	auth.Post("/register", h.AuthHandler.Register)

	protectedAuth := auth.Use(middleware.AuthMiddleware(container.AuthService))
	protectedAuth.Post("/logout", h.AuthHandler.Logout)
	protectedAuth.Post("/logout-all", h.AuthHandler.LogoutAll)
}
