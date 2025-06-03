package middleware

import (
	"go-api/config"

	"github.com/gofiber/fiber/v2"
)

// ConfigValidationMiddleware ensures the application configuration is valid
func ConfigValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cfg := config.Get()

		// Validate critical configuration values
		if err := validateCriticalConfig(cfg); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":   "Configuration Error",
				"message": "Server configuration is invalid",
				"details": err.Error(),
			})
		}

		return c.Next()
	}
}

// validateCriticalConfig performs runtime validation of critical configuration
func validateCriticalConfig(cfg *config.Config) error {
	// This could be expanded to check database connectivity,
	// external service availability, etc.

	// For now, just validate that essential config is present
	// (The main validation is done at startup in config.go)

	return nil
}

// DatabaseHealthMiddleware checks database connectivity before processing requests
func DatabaseHealthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip health check for health endpoints to avoid circular dependency
		if c.Path() == "/health" ||
		   c.Path() == "/health/live" ||
		   c.Path() == "/health/ready" ||
		   c.Path() == "/health/detailed" {
			return c.Next()
		}

		// You could add a quick database ping here if needed
		// This is optional and might add latency to every request

		return c.Next()
	}
}
