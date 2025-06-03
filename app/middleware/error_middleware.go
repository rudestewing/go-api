package middleware

import (
	"go-api/app/shared/errors"
	"go-api/app/shared/response"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// ErrorHandler is a custom error handler middleware for Fiber
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Continue with the next handler
		err := c.Next()
		
		// If there's an error, handle it
		if err != nil {
			return response.HandleAppError(c, err)
		}
		
		return nil
	}
}

// RecoverMiddleware handles panics and converts them to internal server errors
func RecoverMiddleware() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			// Log the panic for debugging
			log.Printf("Panic recovered: %v", e)
			
			// Return internal server error
			response.HandleAppError(c, errors.ErrInternalServer)
		},
	})
}

// NotFoundHandler handles 404 errors
func NotFoundHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return response.HandleAppError(c, errors.NewNotFoundError("Route"))
	}
}

// MethodNotAllowedHandler handles 405 errors
func MethodNotAllowedHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return response.HandleAppError(c, errors.NewAppError(
			"METHOD_NOT_ALLOWED",
			"Method not allowed for this route",
			405,
			nil,
		))
	}
}
