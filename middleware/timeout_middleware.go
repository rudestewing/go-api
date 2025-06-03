package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

// TimeoutMiddleware creates a middleware that enforces a timeout on request processing
// The message parameter is optional - if empty, a default message will be used
func TimeoutMiddleware(timeout time.Duration, message ...string) fiber.Handler {
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
				// Use custom message if provided, otherwise use default
				timeoutMessage := "Request timeout"
				if len(message) > 0 && message[0] != "" {
					timeoutMessage = message[0]
				}
				return fiber.NewError(fiber.StatusRequestTimeout, timeoutMessage)
			}
			return ctx.Err()
		}
	}
}
