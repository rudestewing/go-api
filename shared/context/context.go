package context

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
)

// CreateContext creates a new Go context with values from the Fiber context
func CreateContext(c *fiber.Ctx) context.Context {
	// Start with a background context
	ctx := context.Background()

	// Add any fiber.Locals values you need to pass to the service/repository layer
	if userID := c.Locals("user_id"); userID != nil {
		ctx = context.WithValue(ctx, UserIDKey, userID)
	}

	if email := c.Locals("email"); email != nil {
		ctx = context.WithValue(ctx, EmailKey, email)
	}

	return ctx
}
