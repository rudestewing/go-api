package response

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ValidationResponse represents validation error response structure (Laravel style)
type ValidationResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

// Success sends a successful response
func Success(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *fiber.Ctx, statusCode int, err error, message ...string) error {
	msg := "An error occurred"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	
	response := Response{
		Success: false,
		Message: msg,
	}

	if err != nil {
		response.Error = err.Error()
	}

	return c.Status(statusCode).JSON(response)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *fiber.Ctx, err error, message ...string) error {
	return Error(c, fiber.StatusBadRequest, err, message...)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *fiber.Ctx, message ...string) error {
	return Error(c, fiber.StatusUnauthorized, nil, message...)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *fiber.Ctx, message ...string) error {
	return Error(c, fiber.StatusForbidden, nil, message...)
}

// NotFound sends a 404 Not Found response
func NotFound(c *fiber.Ctx, message ...string) error {
	return Error(c, fiber.StatusNotFound, nil, message...)
}

// UnprocessableEntity sends a 422 Unprocessable Entity response
func UnprocessableEntity(c *fiber.Ctx, err error, message ...string) error {
	return Error(c, fiber.StatusUnprocessableEntity, err, message...)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *fiber.Ctx, err error, message ...string) error {
	return Error(c, fiber.StatusInternalServerError, err, message...)
}

// Created sends a 201 Created response
func Created(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Resource created successfully"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

// NoContent sends a 204 No Content response
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}


// ValidationError sends a Laravel-style validation error response
func ValidationError(c *fiber.Ctx, errors map[string][]string, message ...string) error {
	msg := "Validation failed"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	
	return c.Status(fiber.StatusBadRequest).JSON(ValidationResponse{
		Success: false,
		Message: msg,
		Errors:  errors,
	})
}