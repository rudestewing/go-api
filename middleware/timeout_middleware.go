package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

// TimeoutMiddleware creates a middleware that enforces a timeout on request processing
func TimeoutMiddleware(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()

		// Replace the context in the fiber context
		c.SetUserContext(ctx)

		// Create a channel to receive the result
		done := make(chan error, 1)

		// Run the next handler in a goroutine
		go func() {
			done <- c.Next()
		}()

		// Wait for either completion or timeout
		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				return fiber.NewError(fiber.StatusRequestTimeout, "Request timeout")
			}
			return ctx.Err()
		}
	}
}

// TimeoutMiddlewareWithCustomMessage creates a timeout middleware with a custom timeout message
func TimeoutMiddlewareWithCustomMessage(timeout time.Duration, message string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()

		c.SetUserContext(ctx)

		done := make(chan error, 1)

		go func() {
			done <- c.Next()
		}()

		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				return fiber.NewError(fiber.StatusRequestTimeout, message)
			}
			return ctx.Err()
		}
	}
}
