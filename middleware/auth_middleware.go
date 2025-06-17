package middleware

import (
	authService "go-api/domain/auth/service"
	"go-api/shared/logger"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(service *authService.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
				"code":  "MISSING_AUTH_HEADER",
			})
		}

		// Extract token (must be in "Bearer token" format for security)
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header must be in 'Bearer token' format",
				"code":  "INVALID_AUTH_FORMAT",
			})
		}
		
		token := authHeader[7:] // Remove "Bearer " prefix
		
		// Validate token is not empty
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token cannot be empty",
				"code":  "EMPTY_TOKEN",
			})
		}

		// Get context with timeout from middleware
		ctx := c.UserContext()

		// Validate token using auth service with context
		accessToken, err := service.ValidateToken(ctx, token)
		if err != nil {
			// Log the validation attempt for security monitoring
			logger.Warnf("Token validation failed for IP %s: %v", c.IP(), err)
			
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
				"code":  "INVALID_TOKEN",
			})
		}
		// Validate that the token and user are still valid
		if accessToken == nil || accessToken.User.ID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token or user not found",
				"code":  "TOKEN_USER_NOT_FOUND",
			})
		}

		// Add user info to context for downstream handlers
		c.Locals("user_id", accessToken.UserID)
		c.Locals("user", accessToken.User)
		c.Locals("access_token", accessToken)

		return c.Next()
	}
}
