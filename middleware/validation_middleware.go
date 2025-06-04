package middleware

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// InputValidationMiddleware validates common security issues in input data
func InputValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip validation for GET requests
		if c.Method() == "GET" {
			return c.Next()
		}

		// Get the raw body
		body := c.Body()
		if len(body) == 0 {
			return c.Next()
		}

		// Try to parse as JSON and validate
		var jsonData interface{}
		if err := json.Unmarshal(body, &jsonData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid JSON format",
				"message": "Request body must be valid JSON",
			})
		}

		// Check for potentially dangerous patterns
		bodyStr := string(body)
		if containsSuspiciousPatterns(bodyStr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid input",
				"message": "Request contains potentially malicious content",
			})
		}

		return c.Next()
	}
}

// containsSuspiciousPatterns checks for common attack patterns
func containsSuspiciousPatterns(input string) bool {
	suspiciousPatterns := []string{
		"<script",
		"javascript:",
		"onload=",
		"onerror=",
		"onclick=",
		"eval(",
		"setTimeout(",
		"setInterval(",
		"document.cookie",
		"document.write",
		"window.location",
		"DROP TABLE",
		"DELETE FROM",
		"INSERT INTO",
		"UPDATE SET",
		"UNION SELECT",
		"OR 1=1",
		"' OR '1'='1",
		"../",
		"..\\",
		"cmd.exe",
		"/bin/bash",
		"/bin/sh",
	}

	lowerInput := strings.ToLower(input)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerInput, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// ContentTypeValidationMiddleware ensures proper content types
func ContentTypeValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip validation for GET requests
		if c.Method() == "GET" {
			return c.Next()
		}

		contentType := c.Get("Content-Type")

		// Allow empty content type for requests without body
		if len(c.Body()) == 0 {
			return c.Next()
		}

		// Check if content type is appropriate
		if !strings.HasPrefix(contentType, "application/json") &&
			!strings.HasPrefix(contentType, "multipart/form-data") &&
			!strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
			return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
				"error":   "Unsupported media type",
				"message": "Content-Type must be application/json, multipart/form-data, or application/x-www-form-urlencoded",
			})
		}

		return c.Next()
	}
}
