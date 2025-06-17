package middleware

import (
	"go-api/config"
	"go-api/shared/logger"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// SecurityHeadersMiddleware adds essential security headers
func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prevent XSS attacks
		c.Set("X-XSS-Protection", "1; mode=block")
		
		// Prevent content type sniffing
		c.Set("X-Content-Type-Options", "nosniff")
		
		// Prevent clickjacking
		c.Set("X-Frame-Options", "DENY")
		
		// Enforce HTTPS in production
		cfg := config.Get()
		if cfg.Environment == "production" {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		
		// Prevent information disclosure
		c.Set("X-Powered-By", "")
		c.Set("Server", "")
		
		return c.Next()
	}
}

// InputSanitizationMiddleware sanitizes user input to prevent injection attacks
func InputSanitizationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Sanitize query parameters
		c.Context().QueryArgs().VisitAll(func(key, value []byte) {
			sanitized := sanitizeInput(string(value))
			c.Context().QueryArgs().Set(string(key), sanitized)
		})
		
		// Sanitize headers (but preserve auth headers)
		c.Context().Request.Header.VisitAll(func(key, value []byte) {
			keyStr := string(key)
			if !isAuthHeader(keyStr) {
				sanitized := sanitizeInput(string(value))
				c.Set(keyStr, sanitized)
			}
		})
		
		return c.Next()
	}
}

// sanitizeInput removes potentially dangerous characters and patterns
func sanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	// Remove common SQL injection patterns (basic protection)
	dangerousPatterns := []string{
		"<script", "</script>", "javascript:", "onload=", "onerror=",
		"';", "\";", "/*", "*/", "--", "xp_", "sp_", "exec", "execute",
	}
	
	for _, pattern := range dangerousPatterns {
		input = strings.ReplaceAll(strings.ToLower(input), pattern, "")
	}
	
	return input
}

// isAuthHeader checks if a header is authentication-related and should not be sanitized
func isAuthHeader(header string) bool {	authHeaders := []string{"authorization", "x-api-key", "x-auth-token"}
	headerLower := strings.ToLower(header)
	
	for _, authHeader := range authHeaders {
		if headerLower == authHeader {
			return true
		}
	}
	return false
}

// AdvancedRateLimitMiddleware provides more sophisticated rate limiting
func AdvancedRateLimitMiddleware() fiber.Handler {
	cfg := config.Get()
	
	if !cfg.RateLimitEnabled {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}
	
	return limiter.New(limiter.Config{
		Max:        cfg.RateLimitMax,
		Expiration: cfg.RateLimitWindow,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Rate limit by IP + User-Agent for better protection
			return c.IP() + ":" + c.Get("User-Agent")
		},
		LimitReached: func(c *fiber.Ctx) error {
			logger.Warnf("Rate limit exceeded for IP: %s, User-Agent: %s", 
				c.IP(), c.Get("User-Agent"))
			
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
				"retry_after": cfg.RateLimitWindow.Seconds(),
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		Storage:                nil, // Use in-memory storage
	})
}

// RequestValidationMiddleware validates request size and content
func RequestValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Limit request body size (10MB max)
		const maxBodySize = 10 * 1024 * 1024
		
		if c.Request().Header.ContentLength() > maxBodySize {
			logger.Warnf("Request body too large from IP: %s, Size: %d", 
				c.IP(), c.Request().Header.ContentLength())
			
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error": "Request body too large",
				"max_size": "10MB",
			})
		}
		
		// Validate Content-Type for POST/PUT requests
		if c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "PATCH" {
			contentType := c.Get("Content-Type")
			if contentType == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Content-Type header is required",
				})
			}
			
			validContentTypes := []string{
				"application/json",
				"application/x-www-form-urlencoded",
				"multipart/form-data",
			}
			
			isValid := false
			for _, validType := range validContentTypes {
				if strings.Contains(contentType, validType) {
					isValid = true
					break
				}
			}
			
			if !isValid {
				return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
					"error": "Unsupported Content-Type",
					"supported_types": validContentTypes,
				})
			}
		}
		
		return c.Next()
	}
}
