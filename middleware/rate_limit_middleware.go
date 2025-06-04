package middleware

import (
	"go-api/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware() fiber.Handler {
	cfg := config.Get()

	if !cfg.RateLimitEnabled {
		// Return a no-op middleware if rate limiting is disabled
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	return limiter.New(limiter.Config{
		Max:        cfg.RateLimitMax,
		Expiration: cfg.RateLimitWindow,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Too many requests",
				"message": "Rate limit exceeded. Please try again later.",
			})
		},
	})
}

// AuthRateLimitMiddleware creates a specific rate limiter for authentication endpoints
func AuthRateLimitMiddleware() fiber.Handler {
	cfg := config.Get()

	if !cfg.RateLimitEnabled {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// More restrictive rate limiting for auth endpoints
	return limiter.New(limiter.Config{
		Max:        5,                // 5 login attempts
		Expiration: 15 * time.Minute, // per 15 minutes
		KeyGenerator: func(c *fiber.Ctx) string {
			// Rate limit by IP for login attempts
			return "auth:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Too many login attempts",
				"message": "You have exceeded the maximum number of login attempts. Please try again in 15 minutes.",
			})
		},
	})
}
