package context

import (
	"context"
	"time"

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

// GetUserID gets the user ID from the context
func GetUserID(ctx context.Context) uint {
	value := ctx.Value(UserIDKey)
	if value == nil {
		return 0
	}
	
	// Handle different types that might be stored in the context
	switch v := value.(type) {
	case uint:
		return v
	case int:
		return uint(v)
	case float64:
		return uint(v)
	}
	
	return 0
}

// GetEmail gets the email from the context
func GetEmail(ctx context.Context) string {
	value := ctx.Value(EmailKey)
	if value == nil {
		return ""
	}
	
	// Type assertion with safety check
	if email, ok := value.(string); ok {
		return email
	}
	
	return ""
}

// GenerateContextWithTimeout creates a new context with the specified timeout duration
// and transfers user data from the source context if provided.
// If duration is 0, a default timeout of 30 seconds is used.
func GenerateContextWithTimeout(duration time.Duration, sourceCtx ...context.Context) (context.Context, context.CancelFunc) {
	// Use default timeout of 30 seconds if duration is 0
	if duration == 0 {
		duration = 30 * time.Second
	}

	// Create a new context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	
	// If a source context is provided, copy relevant values from it
	if len(sourceCtx) > 0 && sourceCtx[0] != nil {
		src := sourceCtx[0]
		
		// Transfer user ID if present
		if userID := src.Value(UserIDKey); userID != nil {
			ctx = context.WithValue(ctx, UserIDKey, userID)
		}
		
		// Transfer email if present
		if email := src.Value(EmailKey); email != nil {
			ctx = context.WithValue(ctx, EmailKey, email)
		}
	}
	
	return ctx, cancel
}
