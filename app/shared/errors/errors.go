package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a generic application error
type AppError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	StatusCode int         `json:"-"`
	Details    interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// ValidationError represents validation errors with field-specific messages
type ValidationError struct {
	AppError
	Fields map[string][]string `json:"fields"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation failed: %s", e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(message string, fields map[string][]string) *ValidationError {
	return &ValidationError{
		AppError: AppError{
			Code:       "VALIDATION_ERROR",
			Message:    message,
			StatusCode: http.StatusUnprocessableEntity,
		},
		Fields: fields,
	}
}

// Predefined error instances
var (
	ErrInvalidRequestBody = &AppError{
		Code:       "INVALID_REQUEST_BODY",
		Message:    "Invalid request body format",
		StatusCode: http.StatusBadRequest,
	}

	ErrUserNotFound = &AppError{
		Code:       "USER_NOT_FOUND",
		Message:    "User not found",
		StatusCode: http.StatusNotFound,
	}

	ErrInvalidCredentials = &AppError{
		Code:       "INVALID_CREDENTIALS",
		Message:    "Invalid email or password",
		StatusCode: http.StatusUnauthorized,
	}

	ErrTokenExpired = &AppError{
		Code:       "TOKEN_EXPIRED",
		Message:    "Token has expired",
		StatusCode: http.StatusUnauthorized,
	}

	ErrTokenInvalid = &AppError{
		Code:       "TOKEN_INVALID",
		Message:    "Invalid token",
		StatusCode: http.StatusUnauthorized,
	}

	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "Unauthorized access",
		StatusCode: http.StatusUnauthorized,
	}

	ErrEmailAlreadyExists = &AppError{
		Code:       "EMAIL_ALREADY_EXISTS",
		Message:    "Email address already registered",
		StatusCode: http.StatusConflict,
	}

	ErrInternalServer = &AppError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "Internal server error",
		StatusCode: http.StatusInternalServerError,
	}
)

// NewAppError creates a new application error
func NewAppError(code, message string, statusCode int, details interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    details,
	}
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(message string, details interface{}) *AppError {
	return &AppError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Details:    details,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:       "CONFLICT",
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:       "BAD_REQUEST",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}
