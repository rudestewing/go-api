package response

import (
	"github.com/gofiber/fiber/v2"
)

// Success sends a successful response
func Success(c *fiber.Ctx, data any, message ...string) error {
	msg := "success"
	if len(message) > 0 {
		msg = message[0]
	}

	return c.JSON(fiber.Map{
		"data":    data,
		"message": msg,
	})
}

// Created sends a created response
func Created(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    data,
		"message": "Created successfully",
	})
}

// Error sends an error response
func Error(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"data":    nil,
		"message": message,
	})
}

// ValidationError sends a validation error response
func ValidationError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"data":    nil,
		"message": message,
	})
}

// BadRequest sends a bad request error
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message)
}

// Unauthorized sends an unauthorized error
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message)
}

// InternalError sends an internal server error
func InternalError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, message)
}

// NotFound sends a not found error
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message)
}
