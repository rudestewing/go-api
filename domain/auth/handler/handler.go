package handler

import (
	"go-api/app"
	"go-api/domain/auth/entity"
	"go-api/domain/auth/service"
	"go-api/email"
	"go-api/shared/response"
	"go-api/shared/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	AuthService  *service.AuthService
	EmailService *email.EmailService
}

func NewAuthHandler(p *app.Provider) *AuthHandler {
	return &AuthHandler{
		AuthService:  service.NewAuthService(p),
		EmailService: p.Email,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req entity.LoginRequest
	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, err, "Invalid request body")
	}
	// Validate request
	if validationErrors := validator.ValidateStruct(&req); validationErrors != nil {
		return response.ValidationError(c, validationErrors)
	}

	// Gunakan context dari Fiber yang sudah memiliki timeout dari middleware
	ctx := c.UserContext()
	// Call service with timeout context
	accessToken, err := h.AuthService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.Success(c, fiber.Map{
		"access_token": accessToken.Token,
		"expires_at":   accessToken.ExpiresAt,
		"user":         accessToken.User,
	}, "Login successful")
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req entity.RegisterRequest
	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, err, "Invalid request body")
	}

	// Validate request structure
	if validationErrors := validator.ValidateStruct(&req); validationErrors != nil {
		return response.ValidationError(c, validationErrors)
	}

	// Additional password validation
	if passwordErrors := validator.ValidatePasswordWithDetails(req.Password); passwordErrors != nil {
		return response.ValidationError(c, passwordErrors)
	}

	// Gunakan context dari Fiber yang sudah memiliki timeout dari middleware
	ctx := c.UserContext()
	// Call service with context
	if err := h.AuthService.Register(ctx, &req); err != nil {
		return response.BadRequest(c, err, "Registration failed")
	}

	return response.Success(c, nil, "User registered successfully")
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return response.Unauthorized(c, "Authorization header required")
	}

	// Extract token (assuming Bearer token format)
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	// Use context from Fiber that already has timeout from middleware
	ctx := c.UserContext()

	if err := h.AuthService.Logout(ctx, token); err != nil {
		return response.InternalServerError(c, err, "Failed to logout")
	}

	return response.Success(c, nil, "Logged out successfully")
}

func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	userID := c.Locals("user_id")
	if userID == nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	// Use context from Fiber that already has timeout from middleware
	ctx := c.UserContext()

	if err := h.AuthService.LogoutAll(ctx, userID.(uint)); err != nil {
		return response.InternalServerError(c, err, "Failed to logout from all devices")
	}

	return response.Success(c, nil, "Logged out from all devices successfully")
}
