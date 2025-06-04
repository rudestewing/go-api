package middleware

import (
	"go-api/shared/logger"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandlingMiddleware provides comprehensive error handling with logging
func ErrorHandlingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Continue to next handler
		err := c.Next()

		// If there was an error, handle it
		if err != nil {
			// Get stack trace for debugging
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, false)
			stackTrace := string(buf[:n])
		// Log the error with context
		logger.Errorf("Request failed | IP: %s | Method: %s | Path: %s | Error: %s | Duration: %v | Stack: %s",
			c.IP(), c.Method(), c.Path(), err.Error(), time.Since(start), stackTrace)

			// Determine status code
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Create error response based on status code
			var response fiber.Map
			switch code {
			case fiber.StatusBadRequest:
				response = fiber.Map{
					"error":   "Bad Request",
					"message": "The request could not be understood by the server",
					"code":    code,
				}
			case fiber.StatusUnauthorized:
				response = fiber.Map{
					"error":   "Unauthorized",
					"message": "Authentication is required to access this resource",
					"code":    code,
				}
			case fiber.StatusForbidden:
				response = fiber.Map{
					"error":   "Forbidden",
					"message": "You don't have permission to access this resource",
					"code":    code,
				}
			case fiber.StatusNotFound:
				response = fiber.Map{
					"error":   "Not Found",
					"message": "The requested resource was not found",
					"code":    code,
				}
			case fiber.StatusMethodNotAllowed:
				response = fiber.Map{
					"error":   "Method Not Allowed",
					"message": "The HTTP method is not allowed for this resource",
					"code":    code,
				}
			case fiber.StatusTooManyRequests:
				response = fiber.Map{
					"error":   "Too Many Requests",
					"message": "Rate limit exceeded. Please try again later",
					"code":    code,
				}
			case fiber.StatusUnprocessableEntity:
				response = fiber.Map{
					"error":   "Unprocessable Entity",
					"message": "The request was well-formed but contains semantic errors",
					"code":    code,
				}
			case fiber.StatusInternalServerError:
				response = fiber.Map{
					"error":   "Internal Server Error",
					"message": "An unexpected error occurred. Please try again later",
					"code":    code,
				}
			case fiber.StatusServiceUnavailable:
				response = fiber.Map{
					"error":   "Service Unavailable",
					"message": "The service is temporarily unavailable. Please try again later",
					"code":    code,
				}
			default:
				response = fiber.Map{
					"error":   "Error",
					"message": err.Error(),
					"code":    code,
				}
			}

			return c.Status(code).JSON(response)
		}
		// Log successful requests
		logger.Infof("Request completed | IP: %s | Method: %s | Path: %s | Status: %d | Duration: %v",
			c.IP(), c.Method(), c.Path(), c.Response().StatusCode(), time.Since(start))

		return nil
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate a simple request ID (you might want to use a UUID library)
		requestID := generateRequestID()

		// Set the request ID in the context
		c.Locals("request_id", requestID)

		// Add to response headers
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}

// generateRequestID creates a simple request ID based on timestamp and random component
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
