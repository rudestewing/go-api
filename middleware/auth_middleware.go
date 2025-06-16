package middleware

import (
	authService "go-api/app/domain/auth/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(service *authService.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No authorization header",
			})
		}

		// Extract token (supports both "Bearer token" and just "token")
		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:]
		}

		// Get context with timeout from middleware
		ctx := c.UserContext()

		// Validate token using auth service with context
		accessToken, err := service.ValidateToken(ctx, token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Add user info to context
		c.Locals("user_id", accessToken.UserID)
		c.Locals("user", accessToken.User)
		c.Locals("access_token", accessToken)

		return c.Next()
	}
}
