package router

import (
	"go-api/app"
	"go-api/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(fiberApp *fiber.App, app *app.Provider) {
	h := NewHandler(app)
	
	f := fiberApp.Group("/")
	// HEALTH CHECKS
	f.Get("/health", h.health.HealthCheck)
	f.Get("/health/detailed", h.health.DetailedHealthCheck)
	f.Get("/health/ready", h.health.ReadinessCheck)
	f.Get("/health/live", h.health.LivenessCheck)

	router := fiberApp.Group("/api/v1")
	router.Use(middleware.TimeoutMiddleware(30*time.Second, "Operation timed out"))
	router.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"data":    nil,
			"message": "API is running",
		})
	})	
	
	// AUTHENTICATION ROUTES
	auth := router.Group("/auth")
	auth.Use(middleware.AuthRateLimitMiddleware()) // More restrictive rate limiting for auth	auth.Post("/login", authHandler.Login)
	auth.Post("/register", h.auth.Register)
	auth.Post("/login", h.auth.Login)

	protectedAuth := auth.Use(middleware.AuthMiddleware(app))
	protectedAuth.Post("/logout", h.auth.Logout)
	protectedAuth.Post("/logout-all", h.auth.LogoutAll)
}
