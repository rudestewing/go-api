package handler

import (
	"go-api/app/dto"
	"go-api/app/middleware"
	"go-api/app/service"
	"go-api/app/shared/validator"
	"go-api/container"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService  *service.AuthService
	emailService *service.EmailService
}

func RegisterAuthHandler(router fiber.Router, container *container.Container) { 
	h := &AuthHandler{
		authService:  container.AuthService,
		emailService: container.EmailService,
	}

	public := router.Group("/auth")
	public.Post("/register", h.Register)
	public.Post("/login", h.Login)
	
	protected := router.Use(middleware.AuthMiddleware(container.AuthService))
	protected.Post("/logout", h.Logout)
	protected.Post("/logout-all", h.LogoutAll)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"data":    nil,
			"message": "Invalid request body",
		})
	}

	// Validate request
	if err := validator.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
		})
	}

	// Gunakan context dari Fiber yang sudah memiliki timeout dari middleware
	ctx := c.UserContext()

	// Call service with timeout context
	accessToken, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"access_token": accessToken.Token,
			"expires_at":   accessToken.ExpiresAt,
			"user":         accessToken.User,
		},
		"message": "success",
	})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"data":    nil,
			"message": "Invalid request body",
		})
	}

	// Validate request structure
	if err := validator.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
		})
	}

	// Additional password validation
	if err := validator.ValidatePassword(req.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
		})
	}

	// Gunakan context dari Fiber yang sudah memiliki timeout dari middleware
	ctx := c.UserContext()

	// Call service with context
	if err := h.authService.Register(ctx, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
		})
	}

	// Send welcome email in background (don't block the response)
	go func() {
		if err := h.emailService.SendWelcomeEmail(req.Email, req.Name); err != nil {
			// Log the error but don't fail the registration
			// You might want to use your logger here
			// For now, we'll just continue silently
			// In production, you should log this error properly
		}
	}()

	return c.JSON(fiber.Map{
		"data":    nil,
		"message": "User registered successfully",
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"data":    nil,
			"message": "Authorization header required",
		})
	}

	// Extract token (assuming Bearer token format)
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	// Use context from Fiber that already has timeout from middleware
	ctx := c.UserContext()

	if err := h.authService.Logout(ctx, token); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": "Failed to logout",
		})
	}

	return c.JSON(fiber.Map{
		"data":    nil,
		"message": "Logged out successfully",
	})
}

func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"data":    nil,
			"message": "Unauthorized",
		})
	}

	// Use context from Fiber that already has timeout from middleware
	ctx := c.UserContext()

	if err := h.authService.LogoutAll(ctx, userID.(uint)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": "Failed to logout from all devices",
		})
	}

	return c.JSON(fiber.Map{
		"data":    nil,
		"message": "Logged out from all devices successfully",
	})
}
