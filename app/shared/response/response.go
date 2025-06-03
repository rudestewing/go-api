package response

import (
	"go-api/app/shared/errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// ErrorInfo represents error information in responses
type ErrorInfo struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Fields  map[string][]string    `json:"fields,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Meta represents pagination and other metadata
type Meta struct {
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// Success sends a successful response
func Success(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Operation completed successfully"
	if len(message) > 0 {
		msg = message[0]
	}

	response := APIResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return c.JSON(response)
}

// SuccessWithMeta sends a successful response with metadata
func SuccessWithMeta(c *fiber.Ctx, data interface{}, meta *Meta, message ...string) error {
	msg := "Operation completed successfully"
	if len(message) > 0 {
		msg = message[0]
	}

	response := APIResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return c.JSON(response)
}

// Created sends a created response
func Created(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Resource created successfully"
	if len(message) > 0 {
		msg = message[0]
	}

	response := APIResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// HandleAppError handles application errors and sends appropriate response
func HandleAppError(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case *errors.ValidationError:
		return ValidationErrorResponse(c, e)
	case *errors.AppError:
		return AppErrorResponse(c, e)
	default:
		// For unknown errors, return internal server error
		return InternalServerError(c, "An unexpected error occurred")
	}
}

// ValidationErrorResponse sends a validation error response with field-specific errors
func ValidationErrorResponse(c *fiber.Ctx, validationErr *errors.ValidationError) error {
	response := APIResponse{
		Success: false,
		Message: validationErr.Message,
		Error: &ErrorInfo{
			Code:   validationErr.Code,
			Message: validationErr.Message,
			Fields: validationErr.Fields,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return c.Status(validationErr.StatusCode).JSON(response)
}

// AppErrorResponse sends an application error response
func AppErrorResponse(c *fiber.Ctx, appErr *errors.AppError) error {
	errorDetails := make(map[string]interface{})
	if appErr.Details != nil {
		errorDetails["details"] = appErr.Details
	}

	response := APIResponse{
		Success: false,
		Message: appErr.Message,
		Error: &ErrorInfo{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: errorDetails,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return c.Status(appErr.StatusCode).JSON(response)
}

// Error sends a generic error response
func Error(c *fiber.Ctx, statusCode int, message string) error {
	response := APIResponse{
		Success: false,
		Message: message,
		Error: &ErrorInfo{
			Code:    "ERROR",
			Message: message,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return c.Status(statusCode).JSON(response)
}

// ValidationError sends a validation error response (legacy function)
func ValidationError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnprocessableEntity, message)
}

// BadRequest sends a bad request error
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message)
}

// Unauthorized sends an unauthorized error
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message)
}

// NotFound sends a not found error
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message)
}

// Conflict sends a conflict error
func Conflict(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusConflict, message)
}

// InternalServerError sends an internal server error
func InternalServerError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, message)
}
