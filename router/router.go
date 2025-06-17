package router

import (
	"go-api/container"
	auth "go-api/domain/auth/handler"
	healthcheck "go-api/domain/healthcheck/handler"
	"go-api/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, c *container.Container) {
	// HEALTH CHECKS
	healthHandler := healthcheck.NewHealthHandler(c.DB)
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

	// AUTHENTICATION ROUTES
	authHandler := auth.NewAuthHandler(c.AuthService, c.EmailService)
	auth := router.Group("/auth")
	auth.Use(middleware.AuthRateLimitMiddleware()) // More restrictive rate limiting for auth	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	protectedAuth := auth.Use(middleware.AuthMiddleware(c.AuthService))
	protectedAuth.Post("/logout", authHandler.Logout)
	protectedAuth.Post("/logout-all", authHandler.LogoutAll)
}
